package es.unizar.eina.M27_camping.database;

import androidx.annotation.NonNull;
import androidx.room.ColumnInfo;
import androidx.room.Entity;
import androidx.room.PrimaryKey;

import java.util.Date;

/** Clase aparcelada como entidad que representa una reserva */
@Entity(tableName = "reserva")
public class Reserva {
    @PrimaryKey(autoGenerate = true)
    @ColumnInfo(name = "idReserva")
    private int idReserva;

    @NonNull
    @ColumnInfo(name = "nomCliente")
    private String nomCliente;

    @ColumnInfo(name = "tflCliente")
    private Integer tflCliente;

    @ColumnInfo(name = "fechaEntrada")
    private Date fechaEntrada;

    @ColumnInfo(name = "fechaSalida")
    private Date fechaSalida;

    public Reserva(@NonNull String nomCliente, Integer tflCliente, Date fechaEntrada, Date fechaSalida) {
        this.nomCliente = nomCliente;
        this.tflCliente = tflCliente;
        this.fechaEntrada = fechaEntrada;
        this.fechaSalida = fechaSalida;
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
        return this.tflCliente;
    }

    /** Devuelve la fecha de entrada de la reserva */
    public Date getFechaEntrada() { return this.fechaEntrada; }

    /** Devuelve la fecha de salida de la reserva */
    public Date getFechaSalida() { return this.fechaSalida; }

}
