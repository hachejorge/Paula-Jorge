package es.unizar.eina.M27_camping.database;

import android.content.Context;

import androidx.annotation.NonNull;
import androidx.room.Database;
import androidx.room.Room;
import androidx.room.RoomDatabase;
import androidx.sqlite.db.SupportSQLiteDatabase;

import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

/**
 *   Clase que crea la base de datos de la aplicación ParcelAPP
 *   Principalmente compuesta por parcelas y reservas
 **/
@Database(entities = {Parcela.class, Reserva.class, ParcelaReservada.class}, version = 2, exportSchema = false)
public abstract class CampingRoomDatabase extends RoomDatabase {

    public abstract ParcelaDao parcelaDao();
    public abstract ReservaDao reservaDao();
    public abstract ParcelaReservadaDao parcelaReservadaDao();

    private static volatile CampingRoomDatabase INSTANCE;
    private static final int NUMBER_OF_THREADS = 4;
    static final ExecutorService databaseWriteExecutor =
            Executors.newFixedThreadPool(NUMBER_OF_THREADS);

    static CampingRoomDatabase getDatabase(final Context context) {
        if (INSTANCE == null) {
            synchronized (CampingRoomDatabase.class) {
                if (INSTANCE == null) {
                    INSTANCE = Room.databaseBuilder(context.getApplicationContext(),
                            CampingRoomDatabase.class, "camping_database")
                            .addCallback(sRoomDatabaseCallback)
                            .build();
                }
            }
        }
        return INSTANCE;
    }

    private static RoomDatabase.Callback sRoomDatabaseCallback = new RoomDatabase.Callback() {
        @Override
        public void onCreate(@NonNull SupportSQLiteDatabase db) {
            super.onCreate(db);

            // If you want to keep data through app restarts,
            // comment out the following block
            databaseWriteExecutor.execute(() -> {
                // Populate the database in the background.
                // If you want to start with more notes, just add them.
                ParcelaDao daoParcela = INSTANCE.parcelaDao();
                daoParcela.deleteAll();

                ReservaDao daoReserva = INSTANCE.reservaDao();
                daoReserva.deleteAll();

                Parcela parcela = new Parcela("Alameda", "Bonita parcela de 300m2, es muy acogedora", 8, 15.0f);
                daoParcela.insert(parcela);
                parcela = new Parcela("Pinares", "800m2, una locura", 16, 24.99f);
                daoParcela.insert(parcela);
                parcela = new Parcela("El cerro", "Parela amplia y familiar de 400m2 con una fuente de agua y conexión eléctrica", 10, 12.75f);
                daoParcela.insert(parcela);

                Reserva reserva = new Reserva("Juan José", 620123456, "12/6/2024", "18/6/2024");
                daoReserva.insert(reserva);
                reserva = new Reserva("Pepito", 987654321, "30/12/2024", "2/1/2025");
                daoReserva.insert(reserva);


            });
        }
    };


}
