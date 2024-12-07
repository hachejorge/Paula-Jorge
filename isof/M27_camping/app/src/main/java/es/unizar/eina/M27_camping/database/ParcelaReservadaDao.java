package es.unizar.eina.M27_camping.database;

import androidx.lifecycle.LiveData;
import androidx.room.Dao;
import androidx.room.Delete;
import androidx.room.Insert;
import androidx.room.OnConflictStrategy;
import androidx.room.Query;
import androidx.room.Update;

import java.util.List;

/** Definición de un Data Access Object para las parcelas reservadas */
@Dao
public interface ParcelaReservadaDao {

    @Insert(onConflict = OnConflictStrategy.IGNORE)
    long insert(ParcelaReservada parcelaReservada);

    @Update
    int update(ParcelaReservada parcelaReservada);

    @Delete
    int delete(ParcelaReservada parcelaReservada);

    @Query("DELETE FROM parcelaReservada")
    void deleteAll();

    @Query("SELECT * FROM parcelaReservada")
    LiveData<List<ParcelaReservada>> getAllParcelasReservadas();

    @Query("SELECT * FROM parcelaReservada WHERE idReservaPR = :id_reserva")
    LiveData<List<ParcelaReservada>> getAllParcelasFromReserva(int id_reserva);

}

