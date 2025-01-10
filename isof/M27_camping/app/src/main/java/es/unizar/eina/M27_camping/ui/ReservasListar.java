package es.unizar.eina.M27_camping.ui;

import static androidx.activity.result.contract.ActivityResultContracts.StartActivityForResult;

import android.content.Intent;
import android.os.Bundle;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.AdapterView;
import android.widget.ArrayAdapter;
import android.widget.Spinner;
import android.widget.Toast;

import android.util.Log;

import androidx.activity.result.ActivityResultLauncher;
import androidx.appcompat.app.AlertDialog;
import androidx.appcompat.app.AppCompatActivity;
import androidx.lifecycle.LiveData;
import androidx.lifecycle.ViewModelProvider;
import androidx.recyclerview.widget.LinearLayoutManager;
import androidx.recyclerview.widget.RecyclerView;

import com.google.android.material.floatingactionbutton.FloatingActionButton;

import java.util.List;

import es.unizar.eina.M27_camping.R;
import es.unizar.eina.M27_camping.database.Parcela;
import es.unizar.eina.M27_camping.database.Reserva;
import es.unizar.eina.M27_camping.ui.ReservaEdit;
import es.unizar.eina.send.SendAbstractionImpl;

/** Pantalla en la que se muestra el listado de reservas
 *  Desde ella se pueden añadir, modificar o eliminar cualquiera de ellas. También permite ordenarlas.
 * */
public class ReservasListar extends AppCompatActivity {
    private ReservaViewModel mReservaViewModel;

    static final int INSERT_ID = Menu.FIRST;
    static final int DELETE_ID = Menu.FIRST + 1;
    static final int EDIT_ID = Menu.FIRST + 2;

    private LiveData<List<Reserva>> mReservasOrdenadasPorNombre;
    private LiveData<List<Reserva>> mReservasOrdenadasPorTlf;
    private LiveData<List<Reserva>> mReservasOrdenadasPorFechaEntrada;

    private LiveData<List<Reserva>> currentLiveData;

    RecyclerView mRecyclerView;

    ReservaListAdapter mAdapter;

    FloatingActionButton mFab;

    private SendAbstractionImpl metodoEnvio;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_reservaslistar);
        mRecyclerView = findViewById(R.id.recyclerview_reservas);
        mAdapter = new ReservaListAdapter(new ReservaListAdapter.ReservaDiff());
        mRecyclerView.setAdapter(mAdapter);
        mRecyclerView.setLayoutManager(new LinearLayoutManager(this));

        mReservaViewModel = new ViewModelProvider(this).get(ReservaViewModel.class);

        mReservasOrdenadasPorNombre = mReservaViewModel.getAllReservasPorCliente();
        mReservasOrdenadasPorTlf = mReservaViewModel.getAllReservasPorTlf();
        mReservasOrdenadasPorFechaEntrada = mReservaViewModel.getAllReservasPorFEntrada();

        metodoEnvio = new SendAbstractionImpl(ReservasListar.this,"WhatsApp");

        mFab = findViewById(R.id.fab);
        mFab.setOnClickListener(view -> createReserva());

        // Variables y acciones para el spinner
        Spinner spinnerOrden = findViewById(R.id.spinner);

        String[] opcionesSpinner = {"Nombre", "Teléfono", "Fecha Entrada"};

        ArrayAdapter<String> adapter = new ArrayAdapter<>(this, android.R.layout.simple_spinner_item, opcionesSpinner);
        adapter.setDropDownViewResource(android.R.layout.simple_spinner_dropdown_item);
        spinnerOrden.setAdapter(adapter);

        spinnerOrden.setOnItemSelectedListener(new AdapterView.OnItemSelectedListener() {
            @Override
            public void onItemSelected(AdapterView<?> parent, View view, int position, long id) {
                if (currentLiveData != null) {
                    currentLiveData.removeObservers(ReservasListar.this);
                }

                switch (position) {
                    case 0: // Ordenar por Nombre
                        currentLiveData = mReservasOrdenadasPorNombre;
                        break;
                    case 1: // Ordenar por Ocupantes
                        currentLiveData = mReservasOrdenadasPorTlf;
                        break;
                    case 2: // Ordenar por Precio
                        currentLiveData = mReservasOrdenadasPorFechaEntrada;
                        break;
                }

                currentLiveData.observe(ReservasListar.this, reservas -> {
                    mAdapter.submitList(reservas);
                });
            }

            @Override
            public void onNothingSelected(AdapterView<?> parent) {
                // No hacer nada si no se selecciona nada
            }
        });

        // Asigna el listener para editar
        mAdapter.setOnEditClickListener(reserva -> reservaEdit(reserva));

        mAdapter.setOnDeleteClickListener(reserva -> {
            // Crear un diálogo de alerta para confirmar la eliminación
            new AlertDialog.Builder(this)
                    .setTitle("Confirmar eliminación")
                    .setMessage("¿Estás seguro de que quieres eliminar la reserva #" + reserva.getIdReserva() + "?")
                    .setPositiveButton("Sí", (dialog, which) -> {
                        // Eliminar la parcela si el usuario confirma
                        Toast.makeText(
                                getApplicationContext(),
                                "Borrando Reserva #" + reserva.getIdReserva(),
                                Toast.LENGTH_LONG).show();
                        mReservaViewModel.delete(reserva);
                    })
                    .setNegativeButton("No", (dialog, which) -> {
                        // Cancelar la eliminación
                        dialog.dismiss();
                    })
                    .show();
        });

        mAdapter.setSendClickListener(reserva -> metodoEnvio.send(reserva.getTlfCliente().toString(),
                "Confirmacion Reserva #"+ reserva.getIdReserva() + "\n"
                        + "Estimado/a " + reserva.getNomCliente() + "\n"
                        + "Le confirmamos su reserva con los siguientes detalles:\n"
                        + "Fecha de entrada: " + reserva.getFechaEntrada() + "\n"
                        + "Fecha de salida: " + reserva.getFechaSalida() + "\n"
                        + "Precio total: " + reserva.getPrecioTotal() + "\n"
                        + "Gracias por elegirnos, nos vemos pronto!"));

        // It doesn't affect if we comment the following instruction
        registerForContextMenu(mRecyclerView);

    }

    public boolean onCreateOptionsMenu(Menu menu) {
        boolean result = super.onCreateOptionsMenu(menu);
        menu.add(Menu.NONE, INSERT_ID, Menu.NONE, R.string.add_reserva);
        return result;
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        switch (item.getItemId()) {
            case INSERT_ID:
                createReserva();
                return true;
        }
        return super.onOptionsItemSelected(item);
    }

     public boolean onContextItemSelected(MenuItem item) {
        Reserva current = mAdapter.getCurrent();
        switch (item.getItemId()) {
            case DELETE_ID:
                Toast.makeText(
                        getApplicationContext(),
                        "Borrando " + current.getNomCliente(),
                        Toast.LENGTH_LONG).show();
                mReservaViewModel.delete(current);
                return true;
            case EDIT_ID:
                reservaEdit(current);
                return true;
        }
        return super.onContextItemSelected(item);
    }

    private void createReserva() {
        mStartCreateReserva.launch(new Intent(this, ReservaEdit.class));
    }

    ActivityResultLauncher<Intent> mStartCreateReserva = newActivityResultLauncher(new ExecuteActivityResultReservas() {
        @Override
        public void process(Bundle extras, Reserva reserva) {
            mReservaViewModel.insert(reserva);
        }
    });

    ActivityResultLauncher<Intent> newActivityResultLauncher(ExecuteActivityResultReservas executable) {
        return registerForActivityResult(
                new StartActivityForResult(),
                result -> {
                    if (result.getResultCode() == RESULT_OK) {
                        Bundle extras = result.getData().getExtras();
                        Reserva reserva = new Reserva(extras.getString(ReservaEdit.RESERVA_NOMCLIENTE),
                                                      extras.getInt(ReservaEdit.RESERVA_TLFCLIENTE),
                                                      extras.getString(ReservaEdit.RESERVA_FECHAENTRADA),
                                                      extras.getString(ReservaEdit.RESERVA_FECHASALIDA),
                                                      extras.getFloat(ReservaEdit.RESERVA_PRECIO));
                        executable.process(extras, reserva);
                    }
                });
    }

    private void reservaEdit(Reserva current) {
        Intent intent = new Intent(this, ReservaEdit.class);
        intent.putExtra(ReservaEdit.RESERVA_NOMCLIENTE, current.getNomCliente());
        intent.putExtra(ReservaEdit.RESERVA_TLFCLIENTE, current.getTlfCliente());
        intent.putExtra(ReservaEdit.RESERVA_FECHAENTRADA, current.getFechaEntrada());
        intent.putExtra(ReservaEdit.RESERVA_FECHASALIDA, current.getFechaSalida());
        intent.putExtra(ReservaEdit.RESERVA_ID, current.getIdReserva());
        mStartUpdateReserva.launch(intent);
    }

    ActivityResultLauncher<Intent> mStartUpdateReserva = newActivityResultLauncher(new ExecuteActivityResultReservas() {
        @Override
        public void process(Bundle extras, Reserva reserva) {
            int id = extras.getInt(ReservaEdit.RESERVA_ID);
            reserva.setIdReserva(id);
            mReservaViewModel.update(reserva);
        }
    });

}

interface ExecuteActivityResultReservas {
    void process(Bundle extras, Reserva reserva);
}