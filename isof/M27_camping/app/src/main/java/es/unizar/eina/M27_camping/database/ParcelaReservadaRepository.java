package es.unizar.eina.M27_camping.database;


import android.app.Application;
import android.util.Log;

import androidx.lifecycle.LiveData;

import java.util.List;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.Future;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

/**
 * Clase que gestiona el acceso la fuente de datos de las parcelas reservadas.
 * Interacciona con la base de datos a través de las clases CampingRoomDatabase y ParcelaReservadaDao.
 */
public class ParcelaReservadaRepository {

    private final ParcelaReservadaDao mParcelaReservadaDao;
    private final ParcelaDao mParcelaDao;

    private final long TIMEOUT = 15000;

    /**
     * Constructor de ParcelaReservadaRepository utilizando el contexto de la aplicación para instanciar la base de datos.
     * Alternativamente, se podría estudiar la instanciación del repositorio con una referencia a la base de datos
     * siguiendo el ejemplo de
     * <a href="https://github.com/android/architecture-components-samples/blob/main/BasicSample/app/src/main/java/com/example/android/persistence/DataRepository.java">architecture-components-samples/.../persistence/DataRepository</a>
     */
    public ParcelaReservadaRepository(Application application) {
        CampingRoomDatabase db = CampingRoomDatabase.getDatabase(application);
        mParcelaReservadaDao = db.parcelaReservadaDao();
        mParcelaDao = db.parcelaDao();
    }

    /** Devuelve un objeto de tipo LiveData con todas las parcelas reservadas de todas las reservas.
     * Room ejecuta todas las consultas en un hilo separado.
     * El objeto LiveData notifica a los observadores cuando los datos cambian.
     */
    public LiveData<List<ParcelaReservada>> getAllParcelasReservadas() {
        return mParcelaReservadaDao.getAllParcelasReservadas();
    }

    /** Devuelve un objeto de tipo LiveData con todas las parcelas reservadas de una reserva concreta.
     * Room ejecuta todas las consultas en un hilo separado.
     * El objeto LiveData notifica a los observadores cuando los datos cambian.
     */
    public LiveData<List<ParcelaReservada>> getParcelasFromReserva(int id_reserva) {
        return mParcelaReservadaDao.getAllParcelasFromReserva(id_reserva);
    }

    /**
     * Devuelve una lista con todas las parcelasReservadas que coincidan con el id de una reserva concreta
     */
    public List<ParcelaReservada> getListaParcelasFromReserva(int id_reserva) {
        return mParcelaReservadaDao.getListaParcelasFromReserva(id_reserva);
    }

    /**
     * Actualiza las dependencias de la reserva con los id's correspondientes
     */
    public void actualizarReservasPendientes(){
        mParcelaReservadaDao.actualizarReservasPendientes();
    }

    /**
     *
     */
    public int getMaxOcupParcelaR(int id_parcelaR) {
        return mParcelaReservadaDao.getMaxOcupParcelaR(id_parcelaR);
    }

    /** Inserta una parcela nueva en la base de datos
     * @param parcelaReservada La parcela consta de: un nombre (parcela.getName()) no nulo (parcela.getName()!=null) y no vacío
     *             (parcela.getName().length()>0);, una (parcela.getDescription()) no nula, un máximo de ocupantes, 
     *              (parcela.getMaxOcupantes()) y un precio por persona (parcela.getPrecioPorPersona()).
     * @return Si la parcelaReservada se ha insertado correctamente, devuelve el identificador de la parcela que se ha creado. En caso
     *         contrario, devuelve -1 para indicar el fallo.
     */
    public long insert(ParcelaReservada parcelaReservada) {
        /* Para que la App funcione correctamente y no lance una excepción, la modificación de la
         * base de datos se debe lanzar en un hilo de ejecución separado
         * (databaseWriteExecutor.submit). Para poder sincronizar la recuperación del resultado
         * devuelto por la base de datos, se puede utilizar un Future.
         */
        Future<Long> future = CampingRoomDatabase.databaseWriteExecutor.submit(
                () -> mParcelaReservadaDao.insert(parcelaReservada));
        try {
            return future.get(TIMEOUT, TimeUnit.MILLISECONDS);
        } catch (InterruptedException | ExecutionException | TimeoutException ex) {
            Log.d("ParcelaReservadaRepository", ex.getClass().getSimpleName() + ex.getMessage());
            return -1;
        }
    }

    /** Actualiza una parcela en la base de datos
     * @param parcelaReservada La parcela que se desea actualizar y que consta de: un nombre (parcela.getName()) no nulo (parcela.getName()!=null)
     *              y no vacío (parcela.getName().length()>0);, una (parcela.getDescription()) no nula, un máximo de ocupantes, 
     *              (parcela.getMaxOcupantes()) y un precio por persona (parcela.getPrecioPorPersona()).
     * @return Un valor entero con el número de filas modificadas: 1 si el identificador se corresponde con una parcela
     *         previamente insertada; 0 si no existe previamente una parcela con ese identificador, o hay algún problema
     *         con los atributos.
     */
    public int update(ParcelaReservada parcelaReservada) {
        Future<Integer> future = CampingRoomDatabase.databaseWriteExecutor.submit(
                () -> mParcelaReservadaDao.update(parcelaReservada));
        try {
            return future.get(TIMEOUT, TimeUnit.MILLISECONDS);
        } catch (InterruptedException | ExecutionException | TimeoutException ex) {
            Log.d("ParcelaReservadaRepository", ex.getClass().getSimpleName() + ex.getMessage());
            return -1;
        }
    }


    /** Elimina una parcela en la base de datos.
     * @param parcelaReservada Objeto parcela cuyo atributo identificador (parcela.getId()) contiene la clave primaria de la parcela que se
     *             va a eliminar de la base de datos. Se debe cumplir: parcela.getId() > 0.
     * @return Un valor entero con el número de filas eliminadas: 1 si el identificador se corresponde con una parcela
     *         previamente insertada; 0 si no existe previamente una parcela con ese identificador o el identificador no es
     *         un valor aceptable.
     */
    public int delete(ParcelaReservada parcelaReservada) {
        Future<Integer> future = CampingRoomDatabase.databaseWriteExecutor.submit(
                () -> mParcelaReservadaDao.delete(parcelaReservada));
        try {
            return future.get(TIMEOUT, TimeUnit.MILLISECONDS);
        } catch (InterruptedException | ExecutionException | TimeoutException ex) {
            Log.d("ParcelaReservadaRepository", ex.getClass().getSimpleName() + ex.getMessage());
            return -1;
        }
    }

    public void deleteByParcelaId(int idParcela) {
        Log.d("ParcelaReservadaRepository", "deleteByParcelaId called with id: " + idParcela);
        mParcelaReservadaDao.deleteByParcelaId(idParcela);
    }
}
