package es.unizar.eina.M27_camping.ui;

import static androidx.activity.result.contract.ActivityResultContracts.StartActivityForResult;

import android.content.Intent;
import android.os.Bundle;
import android.view.Menu;
import android.view.MenuItem;
import android.widget.Toast;

import androidx.activity.result.ActivityResultLauncher;
import androidx.appcompat.app.AppCompatActivity;
import androidx.lifecycle.ViewModelProvider;
import androidx.recyclerview.widget.LinearLayoutManager;
import androidx.recyclerview.widget.RecyclerView;

import com.google.android.material.floatingactionbutton.FloatingActionButton;

import es.unizar.eina.M27_camping.R;
import es.unizar.eina.M27_camping.database.Reserva;

/** Pantalla principal de la aplicación Notepad */
public class ReservasListar extends AppCompatActivity {
    private ReservaViewModel mReservaViewModel;

    static final int INSERT_ID = Menu.FIRST;
    static final int DELETE_ID = Menu.FIRST + 1;
    static final int EDIT_ID = Menu.FIRST + 2;

    RecyclerView mRecyclerView;

    ReservaListAdapter mAdapter;

    FloatingActionButton mFab;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_reservaslistar);
        mRecyclerView = findViewById(R.id.recyclerview);
        mAdapter = new ReservaListAdapter(new ReservaListAdapter.ReservaDiff());
        mRecyclerView.setAdapter(mAdapter);
        mRecyclerView.setLayoutManager(new LinearLayoutManager(this));

        mReservaViewModel = new ViewModelProvider(this).get(ReservaViewModel.class);

        mReservaViewModel.getAllReservas().observe(this, reservas -> {
            // Update the cached copy of the notes in the adapter.
            mAdapter.submitList(reservas);
        });

        mFab = findViewById(R.id.fab);
        mFab.setOnClickListener(view -> createReserva());

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

    ActivityResultLauncher<Intent> mStartCreateReserva = newActivityResultLauncher(new ExecuteActivityResult() {
        @Override
        public void process(Bundle extras, Reserva reserva) {
            mReservaViewModel.insert(reserva);
        }
    });

    ActivityResultLauncher<Intent> newActivityResultLauncher(ExecuteActivityResult executable) {
        return registerForActivityResult(
                new StartActivityForResult(),
                result -> {
                    if (result.getResultCode() == RESULT_OK) {
                        Bundle extras = result.getData().getExtras();
                        Reserva reserva = new Reserva(extras.getString(ReservaEdit.RESERVA_NOMCLIENTE),
                                                      extras.getString(ReservaEdit.RESERVA_TLFCLIENTE),
                                                      extras.getInt(ReservaEdit.RESERVA_FECHAENTRADA),
                                                      extras.getFloat(ReservaEdit.RESERVA_FECHASALIDA));
                        executable.process(extras, reserva);
                    }
                });
    }

    private void reservaEditEdit(Reserva current) {
        Intent intent = new Intent(this, ReservaEdit.class);
        intent.putExtra(ReservaEdit.RESERVA_NOMCLIENTE, current.getNomCliente());
        intent.putExtra(ReservaEdit.RESERVA_TLFCLIENTE, current.getTlfCliente());
        intent.putExtra(ReservaEdit.RESERVA_FECHAENTRADA, current.getFechaEntrada());
        intent.putExtra(ReservaEdit.RESERVA_FECHASALIDA, current.getFechaSalida());
        intent.putExtra(ReservaEdit.RESERVA_ID, current.getIdReserva());
        mStartUpdateReserva.launch(intent);
    }

    ActivityResultLauncher<Intent> mStartUpdateReserva = newActivityResultLauncher(new ExecuteActivityResult() {
        @Override
        public void process(Bundle extras, Reserva reserva) {
            int id = extras.getInt(ReservaEdit.RESERVA_ID);
            reserva.setIdReserva(id);
            mReservaViewModel.update(reserva);
        }
    });

}

interface ExecuteActivityResult {
    void process(Bundle extras, Reserva reserva);
}