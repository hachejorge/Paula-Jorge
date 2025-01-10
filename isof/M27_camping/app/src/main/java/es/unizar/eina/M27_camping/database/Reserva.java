package es.unizar.eina.M27_camping.database;

import androidx.annotation.NonNull;
import androidx.room.ColumnInfo;
import androidx.room.Entity;
import androidx.room.PrimaryKey;

import java.text.SimpleDateFormat;
import java.util.Locale;
import java.text.ParseException;
import java.util.Date;

/** Clase Reserva como entidad que representa una reserva formada por un nombre del cliente, un teléfono y una fecha de entrada y de salida */
@Entity(tableName = "reserva")
public class Reserva {
    @PrimaryKey(autoGenerate = true)
    @ColumnInfo(name = "idReserva")
    private int idReserva;

    @NonNull
    @ColumnInfo(name = "nomCliente")
    private String nomCliente;

    @ColumnInfo(name = "tlfCliente")
    private Integer tlfCliente;

    @ColumnInfo(name = "fechaEntrada")
    private String fechaEntrada;

    @ColumnInfo(name = "fechaSalida")
    private String fechaSalida;

    @NonNull
    @ColumnInfo(name = "precioTotal")
    private Float precioTotal;

    private static final SimpleDateFormat dateFormat = new SimpleDateFormat("yyyy-MM-dd", Locale.getDefault());

    public Reserva(@NonNull String nomCliente, Integer tlfCliente, String fechaEntrada, String fechaSalida, Float precioTotal) {
        this.nomCliente = nomCliente;
        this.tlfCliente = tlfCliente;
        this.fechaEntrada = fechaEntrada;
        this.fechaSalida = fechaSalida;
        this.precioTotal = precioTotal;
    }

    /** Devuelve el identificador de la reserva */
    public int getIdReserva(){
        return this.idReserva;
    }

    /** Permite actualizar el identificador de una reserva */
    public void setIdReserva(int idReserva) {
        this.idReserva = idReserva;
    }

    /** Devuelve el nombre del cliente de la reserva */
    public String getNomCliente(){ return this.nomCliente; }

    /** Devuelve el telefono del cliente de la reserva */
    public Integer getTlfCliente(){
        return this.tlfCliente;
    }

    /** Devuelve la fecha de entrada como String */
    public String getFechaEntrada(){
        return this.fechaEntrada;
    }

    /** Devuelve la fecha de salida como String */
    public String getFechaSalida(){
        return this.fechaSalida;
    }

    /** Devuelve el precio total a pagar por la reserva **/
    public Float getPrecioTotal() {return this.precioTotal; }
}
