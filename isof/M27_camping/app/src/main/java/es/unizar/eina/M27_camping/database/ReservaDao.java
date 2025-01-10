package es.unizar.eina.M27_camping.database;

import androidx.lifecycle.LiveData;
import androidx.room.Dao;
import androidx.room.Delete;
import androidx.room.Insert;
import androidx.room.OnConflictStrategy;
import androidx.room.Query;
import androidx.room.Update;

import java.util.List;

/** Definición de un Data Access Object para las reservas */
@Dao
public interface ReservaDao {

    @Insert(onConflict = OnConflictStrategy.IGNORE)
    long insert(Reserva reserva);

    @Update
    int update(Reserva reserva);

    @Delete
    int delete(Reserva reserva);

    @Query("DELETE FROM reserva")
    void deleteAll();

    @Query("SELECT * FROM reserva ORDER BY LOWER(nomCliente) ASC")
    LiveData<List<Reserva>> getOrderedReservasPorCliente();

    @Query("SELECT * FROM reserva ORDER BY tlfCliente ASC")
    LiveData<List<Reserva>> getOrderedReservasPorTlf();

    @Query("SELECT * FROM reserva ORDER BY fechaEntrada ASC")
    LiveData<List<Reserva>> getOrderedReservasPorFechaEntrada();
}

