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

public class UnitTests {

    private static ParcelaDao parcelaDao;
    private static ReservaDao reservaDao;
    private final ExecutorService executor = Executors.newSingleThreadExecutor();

    public UnitTests(CampingRoomDatabase database) {
        this.parcelaDao = database.parcelaDao();
        this.reservaDao = database.reservaDao();
    }



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
