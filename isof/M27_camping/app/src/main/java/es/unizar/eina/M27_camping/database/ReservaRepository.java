package es.unizar.eina.M27_camping.database;


import android.app.Application;
import android.speech.RecognitionService;
import android.util.Log;

import androidx.lifecycle.LiveData;

import java.util.List;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.Future;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

/**
 * Clase que gestiona el acceso la fuente de datos.
 * Interacciona con la base de datos a través de las clases CampingRoomDatabase y ReservaDao.
 */
public class ReservaRepository {

    private final ReservaDao mReservaDao;
    private final LiveData<List<Reserva>> mAllReservas;

    private final long TIMEOUT = 15000;

    /**
     * Constructor de ReservaRepository utilizando el contexto de la aplicación para instanciar la base de datos.
     * Alternativamente, se podría estudiar la instanciación del repositorio con una referencia a la base de datos
     * siguiendo el ejemplo de
     * <a href="https://github.com/android/architecture-components-samples/blob/main/BasicSample/app/src/main/java/com/example/android/persistence/DataRepository.java">architecture-components-samples/.../persistence/DataRepository</a>
     */
    public ReservaRepository(Application application) {
        CampingRoomDatabase db = CampingRoomDatabase.getDatabase(application);
        mReservaDao = db.reservaDao();
        mAllReservas = mReservaDao.getOrderedReservas();
    }

    /** Devuelve un objeto de tipo LiveData con todas las reservas.
     * Room ejecuta todas las consultas en un hilo separado.
     * El objeto LiveData notifica a los observadores cuando los datos cambian.
     */
    public LiveData<List<Reserva>> getAllReservas() {
        return mAllReservas;
    }

    /** Inserta una reserva nueva en la base de datos
     * @param reserva La reserva consta de: un nomCliente (reserva.getNomCliente()) no nulo (reserva.getNomCliente()!=null)
     *      *         y no vacío (reserva.getNomCliente().length()>0);, un (reserva.getTlfCliente()) no nulo, una fecha de entrada,
     *      *         (reserva.getFechaEntrada()) y una fecha de salida (reserva.getFechaSalida()).
     * @return Si la reserva se ha insertado correctamente, devuelve el identificador de la reserva que se ha creado. En caso
     *         contrario, devuelve -1 para indicar el fallo.
     */
    public long insert(Reserva reserva) {
        /* Para que la App funcione correctamente y no lance una excepción, la modificación de la
         * base de datos se debe lanzar en un hilo de ejecución separado
         * (databaseWriteExecutor.submit). Para poder sincronizar la recuperación del resultado
         * devuelto por la base de datos, se puede utilizar un Future.
         */
        Future<Long> future = CampingRoomDatabase.databaseWriteExecutor.submit(
                () -> mReservaDao.insert(reserva));
        try {
            return future.get(TIMEOUT, TimeUnit.MILLISECONDS);
        } catch (InterruptedException | ExecutionException | TimeoutException ex) {
            Log.d("ReservaRepository", ex.getClass().getSimpleName() + ex.getMessage());
            return -1;
        }
    }

    /** Actualiza una reserva en la base de datos
     * @param reserva La reserva que se desea actualizar y que consta de: un nomCliente (reserva.getNomCliente()) no nulo (reserva.getNomCliente()!=null)
     *              y no vacío (reserva.getNomCliente().length()>0);, un (reserva.getTlfCliente()) no nulo, una fecha de entrada,
     *              (reserva.getFechaEntrada()) y una fecha de salida (reserva.getFechaSalida()).
     * @return Un valor entero con el número de filas modificadas: 1 si el identificador se corresponde con una reserva
     *         previamente insertada; 0 si no existe previamente una reserva con ese identificador, o hay algún problema
     *         con los atributos.
     */
    public int update(Reserva reserva) {
        Future<Integer> future = CampingRoomDatabase.databaseWriteExecutor.submit(
                () -> mReservaDao.update(reserva));
        try {
            return future.get(TIMEOUT, TimeUnit.MILLISECONDS);
        } catch (InterruptedException | ExecutionException | TimeoutException ex) {
            Log.d("ReservaRepository", ex.getClass().getSimpleName() + ex.getMessage());
            return -1;
        }
    }


    /** Elimina una reserva en la base de datos.
     * @param reserva Objeto reserva cuyo atributo identificador (reserva.getIdReserva()) contiene la clave primaria de la reserva que se
     *             va a eliminar de la base de datos. Se debe cumplir: reserva.getIdReserva() > 0.
     * @return Un valor entero con el número de filas eliminadas: 1 si el identificador se corresponde con una reserva
     *         previamente insertada; 0 si no existe previamente una reserva con ese identificador o el identificador no es
     *         un valor aceptable.
     */
    public int delete(Reserva reserva) {
        Future<Integer> future = CampingRoomDatabase.databaseWriteExecutor.submit(
                () -> mReservaDao.delete(reserva));
        try {
            return future.get(TIMEOUT, TimeUnit.MILLISECONDS);
        } catch (InterruptedException | ExecutionException | TimeoutException ex) {
            Log.d("ReservaRepository", ex.getClass().getSimpleName() + ex.getMessage());
            return -1;
        }
    }
}
