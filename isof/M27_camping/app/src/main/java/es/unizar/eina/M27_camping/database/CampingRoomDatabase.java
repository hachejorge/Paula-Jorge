package es.unizar.eina.M27_camping.database;

import android.content.Context;

import androidx.annotation.NonNull;
import androidx.room.Database;
import androidx.room.Room;
import androidx.room.RoomDatabase;
import androidx.sqlite.db.SupportSQLiteDatabase;

import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

@Database(entities = {Parcela.class, Reserva.class}, version = 1, exportSchema = false)
public abstract class CampingRoomDatabase extends RoomDatabase {

    public abstract ParcelaDao parcelaDao();

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
                ParcelaDao dao = INSTANCE.parcelaDao();
                dao.deleteAll();

                Parcela parcela = new Parcela("La alameda", "Bonita parcela de 300m2, es muy acogedora", 8, 15.0f);
                dao.insert(parcela);
                parcela = new Parcela("Los pinares", "800m2, una locura", 16, 24.99f);
                dao.insert(parcela);
            });
        }
    };

}
