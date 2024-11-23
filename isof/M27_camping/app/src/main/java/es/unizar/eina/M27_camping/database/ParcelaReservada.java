package es.unizar.eina.M27_camping.database;

import androidx.annotation.NonNull;
import androidx.room.ColumnInfo;
import androidx.room.Entity;
import androidx.room.PrimaryKey;

import java.text.SimpleDateFormat;
import java.util.Locale;

/** Clase aparcelada como entidad que representa una parcela reservada */
@Entity(
        tableName = "parcelaReservada",
        primaryKeys = {"idReservaPR", "idParcelaPR"} //Definimos la clave primaria compuesta
)

public class ParcelaReservada {
    @NonNull
    @ColumnInfo(name = "idReservaPR")
    private int idReservaPR;

    @NonNull
    @ColumnInfo(name = "idParcelaPR")
    private int idParcelaPR;

    @ColumnInfo(name = "nomParcela")
    private String nomParcela;

    @ColumnInfo(name = "numOcupantes")
    private Integer numOcupantes;



    public ParcelaReservada(@NonNull int idReservaPR, @NonNull int idParcelaPR, String nomParcela, Integer numOcupantes) {
        this.idReservaPR = idReservaPR;
        this.idParcelaPR = idParcelaPR;
        this.nomParcela = nomParcela;
        this.numOcupantes = numOcupantes;
    }


    /** Devuelve el identificador de la reserva de la parcela reservada */
    public int getIdReservaPR(){
        return this.idReservaPR;
    }

    /** Permite actualizar el identificador de una reserva de la parcela reservada */
    public void setIdReservaPR(int idReservaPR) {
        this.idReservaPR = idReservaPR;
    }


    /** Devuelve el identificador de la parcela reservada */
    public int getIdParcelaPR(){
        return this.idParcelaPR;
    }

    /** Permite actualizar el identificador de una parcela de la parcela reservada */
    public void setIdParcelaPR(int idParcelaPR) {
        this.idParcelaPR = idParcelaPR;
    }


    /** Devuelve el nombre de la parcela reservada */
    public String getNomParcela(){
        return this.nomParcela;
    }

    /** Permite actualizar el nombre de una parcela reservada */
    public void setNomParcela(String nomParcela) {
        this.nomParcela = nomParcela;
    }


    /** Devuelve el número de ocupantes de la parcela reservada */
    public Integer getNumOcupantes() {
        return numOcupantes;
    }

    /** Permite actualizar el número de ocupantes de la parcela reservada */
    public void setNumOcupantes(Integer numOcupantes) {
        this.numOcupantes = numOcupantes;
    }


}
