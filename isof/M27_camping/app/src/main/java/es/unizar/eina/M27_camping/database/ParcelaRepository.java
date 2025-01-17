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
 * Clase que gestiona el acceso la fuente de datos de las parcelas.
 * Interacciona con la base de datos a través de las clases CampingRoomDatabase y ParcelaDao.
 */
public class ParcelaRepository {

    private final ParcelaDao mParcelaDao;
    //private final LiveData<List<Parcela>> mAllParcelas;

    private final long TIMEOUT = 15000;

    /**
     * Constructor de ParcelaRepository utilizando el contexto de la aplicación para instanciar la base de datos.
     * Alternativamente, se podría estudiar la instanciación del repositorio con una referencia a la base de datos
     * siguiendo el ejemplo de
     * <a href="https://github.com/android/architecture-components-samples/blob/main/BasicSample/app/src/main/java/com/example/android/persistence/DataRepository.java">architecture-components-samples/.../persistence/DataRepository</a>
     */
    public ParcelaRepository(Application application) {
        CampingRoomDatabase db = CampingRoomDatabase.getDatabase(application);
        mParcelaDao = db.parcelaDao();
        //mAllParcelas = mParcelaDao.getOrderedParcelas();
    }

    /** Devuelve un objeto de tipo LiveData con todas las parcelas ordenadas por Nombre.
     * Room ejecuta todas las consultas en un hilo separado.
     * El objeto LiveData notifica a los observadores cuando los datos cambian.
     */
    public LiveData<List<Parcela>> getAllParcelasPorNombre() {
        return mParcelaDao.getOrderedParcelasByNombre();
    }

    /**
     * Devuelve un objeto de tipo LiveData que contiene la parcela con el nombre especificado.
     * Room ejecuta todas las consultas en un hilo separado.
     * El objeto LiveData notifica a los observadores cuando los datos cambian.
     *
     * @param nom El nombre de la parcela a buscar.
     * @return Un objeto LiveData que contiene la parcela con el nombre especificado, si existe.
     */
    public LiveData<Parcela> getParcelaPorNombre(String nom) {
        return mParcelaDao.getParcelaPorNombre(nom);
    }

    /**
     * Devuelve una parcela de forma síncrona, buscando por su identificador único.
     *
     * @param id El identificador único de la parcela.
     * @return La parcela correspondiente al identificador proporcionado, si existe. Si no existe, devuelve null.
     */
    public Parcela getParcelaPorId(int id) {
        return mParcelaDao.getParcelaById(id);
    }

    /** Devuelve un objeto de tipo LiveData con todas las parcelas ordenadas por Número de ocupantes máximos.
     * Room ejecuta todas las consultas en un hilo separado.
     * El objeto LiveData notifica a los observadores cuando los datos cambian.
     */
    public LiveData<List<Parcela>> getAllParcelasPorOcupantes() {
        return mParcelaDao.getOrderedParcelasByOcupantes();
    }

    /** Devuelve un objeto de tipo LiveData con todas las parcelas ordenadas por precio por persona.
     * Room ejecuta todas las consultas en un hilo separado.
     * El objeto LiveData notifica a los observadores cuando los datos cambian.
     */
    public LiveData<List<Parcela>> getAllParcelasPorPrecio() {
        return mParcelaDao.getOrderedParcelasByPrecio();
    }

    /** Inserta una parcela nueva en la base de datos
     * @param parcela La parcela consta de: un nombre (parcela.getName()) no nulo (parcela.getName()!=null) y no vacío
     *             (parcela.getName().length()>0);, una (parcela.getDescription()) no nula, un máximo de ocupantes, 
     *              (parcela.getMaxOcupantes()) y un precio por persona (parcela.getPrecioPorPersona()).
     * @return Si la parcela se ha insertado correctamente, devuelve el identificador de la parcela que se ha creado. En caso
     *         contrario, devuelve -1 para indicar el fallo.
     */
    public long insert(Parcela parcela) {
        /* Para que la App funcione correctamente y no lance una excepción, la modificación de la
         * base de datos se debe lanzar en un hilo de ejecución separado
         * (databaseWriteExecutor.submit). Para poder sincronizar la recuperación del resultado
         * devuelto por la base de datos, se puede utilizar un Future.
         */
        Future<Long> future = CampingRoomDatabase.databaseWriteExecutor.submit(
                () -> mParcelaDao.insert(parcela));
        try {
            return future.get(TIMEOUT, TimeUnit.MILLISECONDS);
        } catch (InterruptedException | ExecutionException | TimeoutException ex) {
            Log.d("ParcelaRepository", ex.getClass().getSimpleName() + ex.getMessage());
            return -1;
        }
    }

    /** Actualiza una parcela en la base de datos
     * @param parcela La parcela que se desea actualizar y que consta de: un nombre (parcela.getName()) no nulo (parcela.getName()!=null) 
     *              y no vacío (parcela.getName().length()>0);, una (parcela.getDescription()) no nula, un máximo de ocupantes, 
     *              (parcela.getMaxOcupantes()) y un precio por persona (parcela.getPrecioPorPersona()).
     * @return Un valor entero con el número de filas modificadas: 1 si el identificador se corresponde con una parcela
     *         previamente insertada; 0 si no existe previamente una parcela con ese identificador, o hay algún problema
     *         con los atributos.
     */
    public int update(Parcela parcela) {
        Future<Integer> future = CampingRoomDatabase.databaseWriteExecutor.submit(
                () -> mParcelaDao.update(parcela));
        try {
            return future.get(TIMEOUT, TimeUnit.MILLISECONDS);
        } catch (InterruptedException | ExecutionException | TimeoutException ex) {
            Log.d("ParcelaRepository", ex.getClass().getSimpleName() + ex.getMessage());
            return -1;
        }
    }


    /** Elimina una parcela en la base de datos.
     * @param parcela Objeto parcela cuyo atributo identificador (parcela.getId()) contiene la clave primaria de la parcela que se
     *             va a eliminar de la base de datos. Se debe cumplir: parcela.getId() > 0.
     * @return Un valor entero con el número de filas eliminadas: 1 si el identificador se corresponde con una parcela
     *         previamente insertada; 0 si no existe previamente una parcela con ese identificador o el identificador no es
     *         un valor aceptable.
     */
    public int delete(Parcela parcela) {
        Future<Integer> future = CampingRoomDatabase.databaseWriteExecutor.submit(
                () -> mParcelaDao.delete(parcela));
        try {
            return future.get(TIMEOUT, TimeUnit.MILLISECONDS);
        } catch (InterruptedException | ExecutionException | TimeoutException ex) {
            Log.d("ParcelaRepository", ex.getClass().getSimpleName() + ex.getMessage());
            return -1;
        }
    }
}
