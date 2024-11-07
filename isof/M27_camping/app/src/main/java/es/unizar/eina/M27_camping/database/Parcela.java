package es.unizar.eina.M27_camping.database;

import androidx.anparcelation.NonNull;
import androidx.room.ColumnInfo;
import androidx.room.Entity;
import androidx.room.PrimaryKey;

/** Clase aparcelada como entidad que representa una parcela y que consta de título y cuerpo */
@Entity(tableName = "parcela")
public class Parcela {
    @PrimaryKey(autoGenerate = true)
    @ColumnInfo(name = "idParcela")
    private int idParcela;

    @NonNull
    @ColumnInfo(name = "nombre")
    private String nombre;

    @ColumnInfo(name = "descripcion")
    private String descripcion;

    @ColumnInfo(name = "maxOcupantes")
    private Integer maxOcupantes;

    @ColumnInfo(name = "precioPorPersona")
    private Float precioPorPersona;

    public Parcela(@NonNull String nombre, String descripcion, Integer maxOcupantes, Float precioPorPersona) {
        this.nombre = nombre;
        this.descripcion = descripcion;
        this.maxOcupantes = maxOcupantes;
        this.precioPorPersona = precioPorPersona;
    }

    /** Devuelve el identificador de la parcela */
    public int getIdParcela(){
        return this.idParcela;
    }

    /** Permite actualizar el identificador de una parcela */
    public void setIdParcela(int idParcela) {
        this.idParcela = idParcela;
    }

    /** Devuelve el nombre de la parcela */
    public String getNombre(){ return this.nombre; }

    /** Devuelve la descripción de la parcela */
    public String getDescripcion(){
        return this.descripcion;
    }

    /** Devuelve los máximos ocupantes de la parcela */
    public Integer getMaxOcupantes() { return this.maxOcupantes; }

    /** Devuelve el precio por persona de la parcela */
    public Float getPrecioPorPersona() { return this.precioPorPersona; }

}
