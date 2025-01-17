package es.unizar.eina.M27_camping.ui;

import static android.content.ContentValues.TAG;

import android.util.Log;

import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

import es.unizar.eina.M27_camping.database.CampingRoomDatabase;
import es.unizar.eina.M27_camping.database.Parcela;
import es.unizar.eina.M27_camping.database.ParcelaDao;
import es.unizar.eina.M27_camping.database.ParcelaReservadaDao;
import es.unizar.eina.M27_camping.database.Reserva;
import es.unizar.eina.M27_camping.database.ReservaDao;


/**
 * Clase UnitTests para realizar pruebas de funcionalidad sobre la base de datos de la aplicación.
 */
public class UnitTests {

    private static ParcelaDao parcelaDao;
    private static ReservaDao reservaDao;
    private final ExecutorService executor = Executors.newSingleThreadExecutor();

    /**
     * Constructor de la clase UnitTests.
     *
     * @param database Instancia de la base de datos utilizada para obtener los DAOs.
     */
    public UnitTests(CampingRoomDatabase database) {
        this.parcelaDao = database.parcelaDao();
        this.reservaDao = database.reservaDao();
    }


    /**
     * Inserta datos válidos en la base de datos para realizar pruebas.
     */
    public void insertValidos() {
        //Parcelas
        Parcela parcela = new Parcela("El corral", "Muy amplia", 5, 10.0f);
        executor.execute(() -> parcelaDao.insert(parcela));
        //Reservas
        Reserva reserva = new Reserva("Ana", 123456789, "2025-02-20", "2025-02-25", 164.97f * 6);
        executor.execute(() -> reservaDao.insert(reserva));
    }

    /**
     * Inserta datos no válidos en la base de datos para realizar pruebas de validación.
     */
    public void insertNoValidos() {
        //Parcelas
        Parcela parcela = new Parcela(null, "Muy amplia", 5, 10.0f);
        executor.execute(() -> parcelaDao.insert(parcela));
        Parcela parcela1 = new Parcela("", "Muy amplia", 5, 10.0f);
        executor.execute(() -> parcelaDao.insert(parcela1));
        Parcela parcela2 = new Parcela("El corral",null, 5, 10.0f);
        executor.execute(() -> parcelaDao.insert(parcela2));
        Parcela parcela3 = new Parcela("El corral", "", 5, 10.0f);
        executor.execute(() -> parcelaDao.insert(parcela3));
        Parcela parcela4 = new Parcela("El corral", "Muy amplia",null, 10.0f);
        executor.execute(() -> parcelaDao.insert(parcela4));
        Parcela parcela6 = new Parcela("El corral", "Muy amplia", -2, 10.0f);
        executor.execute(() -> parcelaDao.insert(parcela6));
        Parcela parcela7 = new Parcela("El corral", "Muy amplia", 5, null);
        executor.execute(() -> parcelaDao.insert(parcela7));
        Parcela parcela9 = new Parcela("El corral", "Muy amplia", 5, -10.0f);
        executor.execute(() -> parcelaDao.insert(parcela9));
        //Reservas
        Reserva reserva = new Reserva(null, 123456789, "2025-02-20", "2025-02-25", 164.97f * 6);
        executor.execute(() -> reservaDao.insert(reserva));
        Reserva reserva1 = new Reserva("", 123456789, "2025-02-20", "2025-02-25", 164.97f * 6);
        executor.execute(() -> reservaDao.insert(reserva1));
        Reserva reserva2 = new Reserva("Ana", null, "2025-02-20", "2025-02-25", 164.97f * 6);
        executor.execute(() -> reservaDao.insert(reserva2));
        Reserva reserva3 = new Reserva("Ana", 123, "2025-02-20", "2025-02-25", 164.97f * 6);
        executor.execute(() -> reservaDao.insert(reserva3));
        Reserva reserva4 = new Reserva("Ana", 123456789, null, "2025-02-25", 164.97f * 6);
        executor.execute(() -> reservaDao.insert(reserva4));
        Reserva reserva5 = new Reserva("Ana", 123456789, "", "2025-02-25", 164.97f * 6);
        executor.execute(() -> reservaDao.insert(reserva5));
        Reserva reserva6 = new Reserva("Ana", 123456789, "20-02-2025", "2025-02-25", 164.97f * 6);
        executor.execute(() -> reservaDao.insert(reserva6));
        Reserva reserva7 = new Reserva("Ana", 123456789, "1300-34-35", "2025-02-25", 164.97f * 6);
        executor.execute(() -> reservaDao.insert(reserva7));
        Reserva reserva8 = new Reserva("Ana", 123456789, "2025-02-20", null, 164.97f * 6);
        executor.execute(() -> reservaDao.insert(reserva8));
        Reserva reserva9 = new Reserva("Ana", 123456789, "2025-02-20", "", 164.97f * 6);
        executor.execute(() -> reservaDao.insert(reserva9));
        Reserva reserva10 = new Reserva("Ana", 123456789, "2025-02-20", "25-02-2025", 164.97f * 6);
        executor.execute(() -> reservaDao.insert(reserva10));
        Reserva reserva11 = new Reserva("Ana", 123456789, "2025-02-20", "1300-34-35", 164.97f * 6);
        executor.execute(() -> reservaDao.insert(reserva11));
        Reserva reserva12 = new Reserva("Ana", 123456789, "2025-02-20", "25-02-2025",null);
        executor.execute(() -> reservaDao.insert(reserva12));
        Reserva reserva13 = new Reserva("Ana", 123456789, "2025-02-20", "25-02-2025",-6.0f);
        executor.execute(() -> reservaDao.insert(reserva13));
    }


    /**
     * Actualiza datos válidos en la base de datos para realizar pruebas.
     */
    public void updateValidos() {
        //Parcelas
        Parcela parcela = new Parcela("El corral", "Muy amplia", 5, 10.0f);
        executor.execute(() -> parcelaDao.update(parcela));
        //Reservas
        Reserva reserva = new Reserva("Ana", 123456789, "2025-02-20", "2025-02-25", 164.97f * 6);
        executor.execute(() -> reservaDao.update(reserva));
    }


    /**
     * Actualiza datos no válidos en la base de datos para realizar pruebas de validación.
     */
    public void updateNoValidos() {
        //Parcelas
        Parcela parcela = new Parcela(null, "Muy amplia", 5, 10.0f);
        executor.execute(() -> parcelaDao.update(parcela));
        Parcela parcela1 = new Parcela("", "Muy amplia", 5, 10.0f);
        executor.execute(() -> parcelaDao.update(parcela1));
        Parcela parcela2 = new Parcela("El corral",null, 5, 10.0f);
        executor.execute(() -> parcelaDao.update(parcela2));
        Parcela parcela3 = new Parcela("El corral", "", 5, 10.0f);
        executor.execute(() -> parcelaDao.update(parcela3));
        Parcela parcela4 = new Parcela("El corral", "Muy amplia",null, 10.0f);
        executor.execute(() -> parcelaDao.update(parcela4));
        Parcela parcela6 = new Parcela("El corral", "Muy amplia", -2, 10.0f);
        executor.execute(() -> parcelaDao.update(parcela6));
        Parcela parcela7 = new Parcela("El corral", "Muy amplia", 5, null);
        executor.execute(() -> parcelaDao.update(parcela7));
        Parcela parcela9 = new Parcela("El corral", "Muy amplia", 5, -10.0f);
        executor.execute(() -> parcelaDao.update(parcela9));
        //Reservas
        Reserva reserva = new Reserva(null, 123456789, "2025-02-20", "2025-02-25", 164.97f * 6);
        executor.execute(() -> reservaDao.update(reserva));
        Reserva reserva1 = new Reserva("", 123456789, "2025-02-20", "2025-02-25", 164.97f * 6);
        executor.execute(() -> reservaDao.update(reserva1));
        Reserva reserva2 = new Reserva("Ana", null, "2025-02-20", "2025-02-25", 164.97f * 6);
        executor.execute(() -> reservaDao.update(reserva2));
        Reserva reserva3 = new Reserva("Ana", 123, "2025-02-20", "2025-02-25", 164.97f * 6);
        executor.execute(() -> reservaDao.update(reserva3));
        Reserva reserva4 = new Reserva("Ana", 123456789, null, "2025-02-25", 164.97f * 6);
        executor.execute(() -> reservaDao.update(reserva4));
        Reserva reserva5 = new Reserva("Ana", 123456789, "", "2025-02-25", 164.97f * 6);
        executor.execute(() -> reservaDao.update(reserva5));
        Reserva reserva6 = new Reserva("Ana", 123456789, "20-02-2025", "2025-02-25", 164.97f * 6);
        executor.execute(() -> reservaDao.update(reserva6));
        Reserva reserva7 = new Reserva("Ana", 123456789, "1300-34-35", "2025-02-25", 164.97f * 6);
        executor.execute(() -> reservaDao.update(reserva7));
        Reserva reserva8 = new Reserva("Ana", 123456789, "2025-02-20", null, 164.97f * 6);
        executor.execute(() -> reservaDao.update(reserva8));
        Reserva reserva9 = new Reserva("Ana", 123456789, "2025-02-20", "", 164.97f * 6);
        executor.execute(() -> reservaDao.update(reserva9));
        Reserva reserva10 = new Reserva("Ana", 123456789, "2025-02-20", "25-02-2025", 164.97f * 6);
        executor.execute(() -> reservaDao.update(reserva10));
        Reserva reserva11 = new Reserva("Ana", 123456789, "2025-02-20", "1300-34-35", 164.97f * 6);
        executor.execute(() -> reservaDao.update(reserva11));
        Reserva reserva12 = new Reserva("Ana", 123456789, "2025-02-20", "25-02-2025",null);
        executor.execute(() -> reservaDao.update(reserva12));
        Reserva reserva13 = new Reserva("Ana", 123456789, "2025-02-20", "25-02-2025",-6.0f);
        executor.execute(() -> reservaDao.update(reserva13));
    }


    /**
     * Elimina datos válidos de la base de datos para realizar pruebas.
     */
    public void deleteValidos() {
        //Parcelas
        executor.execute(() -> parcelaDao.delete(new Parcela("nombre", "descripcion", 10, 20.0f)));
        //Reservas
        Reserva reserva = new Reserva("cliente", 123456789, "2025-02-20", "2025-02-25", 164.97f * 6);
        executor.execute(() -> reservaDao.delete(reserva));
    }


    /**
     * Elimina datos no válidos de la base de datos para realizar pruebas.
     */
    public void deleteNoValidos() {
        //Parcelas
        executor.execute(() -> parcelaDao.delete(new Parcela(null, "descripcion", 10, 15.0f)));
        //Reservas
        Reserva reserva = new Reserva(null, 123456789, "2025-02-20", "2025-02-25", 164.97f * 6);
        executor.execute(() -> reservaDao.delete(reserva));
    }


    /**
     * Realiza una prueba de caja negra, ejecutando operaciones de inserción, actualización y eliminación.
     */
    public void pruebaCajaNegra() {
        insertValidos();
        insertNoValidos();
        updateValidos();
        updateNoValidos();
        deleteValidos();
        deleteNoValidos();
    }

    /**
     * Realiza una prueba de volumen insertando grandes cantidades de datos en la base de datos.
     */
    public void pruebaVolumen() {
        for(int i=0; i<100; i++) {
            Parcela parcela = new Parcela("Parcela" + String.valueOf(i) , "Parcela prueba " + String.valueOf(i), 5, 10.0f);
            executor.execute(() -> parcelaDao.insert(parcela));
        }

        for(int i=0; i<10000; i++) {
            Reserva reserva = new Reserva("Cliente" + String.valueOf(i), 123456789, "2024-06-12", "2024-06-18", 164.97f * 6);
            executor.execute(() -> reservaDao.insert(reserva));
        }
    }

    /**
     * Realiza una prueba de sobrecarga incrementando gradualmente la longitud de un campo.
     */
    public void pruebaSobrecarga() {
        int longitud = 1000; // Longitud inicial
        StringBuilder descripcionBuilder = new StringBuilder("X".repeat(longitud));

        while (true) {
            String descripcion = descripcionBuilder.toString();
            Parcela parcela = new Parcela("TestParcela", descripcion, 5, 10.0f);
            try {
                executor.execute(() -> {
                    parcelaDao.insert(parcela);
                    Log.d(TAG, "Inserción exitosa con longitud: " + descripcion.length());
                });

                // Incrementar la longitud de la descripción
                longitud += 1000;
                descripcionBuilder.append("X".repeat(100));

            } catch (Exception e) {
                Log.d(TAG, "Fallo al insertar parcela con longitud: " + descripcion.length());
                Log.d(TAG, "Excepción: " + e.getMessage());
                break;
            }
        }
    }


}
