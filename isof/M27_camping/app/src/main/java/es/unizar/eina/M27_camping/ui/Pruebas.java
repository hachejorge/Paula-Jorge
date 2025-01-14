package es.unizar.eina.M27_camping.ui;

import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.widget.Toast;

import androidx.appcompat.app.AppCompatActivity;

import es.unizar.eina.M27_camping.R;
import es.unizar.eina.M27_camping.database.CampingRoomDatabase;
import es.unizar.eina.M27_camping.ui.UnitTests;

public class Pruebas extends AppCompatActivity {

    private UnitTests mUnitTests;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_pruebas);

        CampingRoomDatabase database = CampingRoomDatabase.getDatabase(this);
        mUnitTests = new UnitTests(database);



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
