package es.unizar.eina.M27_camping.ui;

import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.widget.Toast;

import androidx.appcompat.app.AppCompatActivity;

import es.unizar.eina.M27_camping.R;
import es.unizar.eina.M27_camping.database.CampingRoomDatabase;
import es.unizar.eina.M27_camping.ui.UnitTests;

/**
 * Clase que gestiona la interfaz de usuario para ejecutar pruebas del sistema.
 * Proporciona botones para realizar pruebas de caja negra, pruebas de volumen y pruebas de sobrecarga.
 */
public class Pruebas extends AppCompatActivity {

    /**
     * Objeto de la clase UnitTests que contiene las pruebas definidas.
     */
    private UnitTests mUnitTests;

    /**
     * Método llamado al crear la actividad.
     * Inicializa los elementos de la interfaz de usuario y configura los listeners para los botones de prueba.
     *
     * @param savedInstanceState Contiene el estado previamente guardado de la actividad, si existe.
     */
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_pruebas);

        CampingRoomDatabase database = CampingRoomDatabase.getDatabase(this);
        mUnitTests = new UnitTests(database);

        findViewById(R.id.button_pruebas_cajanegra).setOnClickListener( view -> {
            try {
                mUnitTests.pruebaCajaNegra();
                Toast.makeText(this, "Prueba de caja negra realizada con éxito", Toast.LENGTH_SHORT).show();
            } catch (Exception e) {
                Toast.makeText(this, "Error en la prueba de caja negra: " + e.getMessage(), Toast.LENGTH_LONG).show();
            }
        });


        findViewById(R.id.button_prueba_volumen).setOnClickListener( view -> {
            try {
                mUnitTests.pruebaVolumen();
                Toast.makeText(this, "Prueba de volumen realizada con éxito", Toast.LENGTH_SHORT).show();
            } catch (Exception e) {
                Toast.makeText(this, "Error en la prueba de volumen: " + e.getMessage(), Toast.LENGTH_LONG).show();
            }
        });

        findViewById(R.id.button_prueba_sobrecarga).setOnClickListener( view -> {
            try {
                mUnitTests.pruebaSobrecarga();
                Toast.makeText(this, "Prueba de sobrecarga realizada con éxito", Toast.LENGTH_SHORT).show();
            } catch (Exception e) {
                Toast.makeText(this, "Error en la prueba de sobrecarga: " + e.getMessage(), Toast.LENGTH_LONG).show();
            }
        });
    }


}
