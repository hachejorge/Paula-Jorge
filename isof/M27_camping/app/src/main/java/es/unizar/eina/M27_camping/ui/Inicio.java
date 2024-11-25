package es.unizar.eina.M27_camping.ui;

import android.content.Intent;
import android.os.Bundle;
import android.view.View;

import androidx.appcompat.app.AppCompatActivity;
import es.unizar.eina.M27_camping.R;

/**
 * Pantalla inicial de la aplicación ParcelApp
 * En ella te permite acceder a ĺos listados de parcelas y reservas
 * */
public class Inicio extends AppCompatActivity {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        // Se asigna la vista de activity de inicio
        setContentView(R.layout.activity_inicio);

        // Configuramos el botón de ver parcelas
        findViewById(R.id.button_parcelas_listar).setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View view) {
                Intent intent = new Intent(Inicio.this, ParcelaReservadaListar.class);
                startActivity(intent);
            }
        });

        // Configuramos el botón de ver reservas
        findViewById(R.id.button_reservas_listar).setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View view) {
                // Puedes dejar esto como un Toast temporal hasta que implementes ReservasListar
                // startActivity(new Intent(MainActivity.this, ReservasListar.class));
                //Toast.makeText(Inicio.this, "ReservasListar aún no está implementado", Toast.LENGTH_SHORT).show();

                Intent intent = new Intent(Inicio.this, ReservasListar.class);
                startActivity(intent);
            }
        });
    }
}
