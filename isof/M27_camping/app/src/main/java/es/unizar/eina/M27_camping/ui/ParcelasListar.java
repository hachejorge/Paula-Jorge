package es.unizar.eina.M27_camping.ui;

import androidx.activity.result.ActivityResultLauncher;
import androidx.appcompat.app.AlertDialog;
import androidx.appcompat.app.AppCompatActivity;
import androidx.lifecycle.ViewModelProvider;
import androidx.recyclerview.widget.LinearLayoutManager;
import androidx.recyclerview.widget.RecyclerView;

import android.content.Intent;
import android.os.Bundle;
import android.view.Menu;
import android.view.MenuItem;
import android.widget.Toast;

import com.google.android.material.floatingactionbutton.FloatingActionButton;

import es.unizar.eina.M27_camping.database.Parcela;
import es.unizar.eina.M27_camping.R;

import static androidx.activity.result.contract.ActivityResultContracts.StartActivityForResult;

/** Pantalla principal de la aplicación ParcelApp */
public class ParcelasListar extends AppCompatActivity {
    private ParcelaViewModel mParcelaViewModel;

    static final int INSERT_ID = Menu.FIRST;
    static final int DELETE_ID = Menu.FIRST + 1;
    static final int EDIT_ID = Menu.FIRST + 2;

    RecyclerView mRecyclerView;

    ParcelaListAdapter mAdapter;

    FloatingActionButton mFab;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_parcelaslistar);
        mRecyclerView = findViewById(R.id.recyclerview);
        mAdapter = new ParcelaListAdapter(new ParcelaListAdapter.ParcelaDiff());
        mRecyclerView.setAdapter(mAdapter);
        mRecyclerView.setLayoutManager(new LinearLayoutManager(this));

        mParcelaViewModel = new ViewModelProvider(this).get(ParcelaViewModel.class);

        mParcelaViewModel.getAllParcelas().observe(this, parcelas -> {
            // Update the cached copy of the notes in the adapter.
            mAdapter.submitList(parcelas);
        });

        mFab = findViewById(R.id.fab);
        mFab.setOnClickListener(view -> createParela());

        // Asigna el listener para editar
        mAdapter.setOnEditClickListener(parcela -> parelaEdit(parcela));

        mAdapter.setOnDeleteClickListener(parcela -> {
            // Crear un diálogo de alerta para confirmar la eliminación
            new AlertDialog.Builder(this)
                    .setTitle("Confirmar eliminación")
                    .setMessage("¿Estás seguro de que quieres eliminar " + parcela.getNombre() + "?")
                    .setPositiveButton("Sí", (dialog, which) -> {
                        // Eliminar la parcela si el usuario confirma
                        Toast.makeText(
                                getApplicationContext(),
                                "Borrando " + parcela.getNombre(),
                                Toast.LENGTH_LONG).show();
                        mParcelaViewModel.delete(parcela);
                    })
                    .setNegativeButton("No", (dialog, which) -> {
                        // Cancelar la eliminación
                        dialog.dismiss();
                    })
                    .show();
        });


        // It doesn't affect if we comment the following instruction
        // registerForContextMenu(mRecyclerView);

    }

    public boolean onCreateOptionsMenu(Menu menu) {
        boolean result = super.onCreateOptionsMenu(menu);
        menu.add(Menu.NONE, INSERT_ID, Menu.NONE, R.string.add_parcela);
        return result;
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        switch (item.getItemId()) {
            case INSERT_ID:
                createParela();
                return true;
        }
        return super.onOptionsItemSelected(item);
    }

     public boolean onContextItemSelected(MenuItem item) {
        Parcela current = mAdapter.getCurrent();
        switch (item.getItemId()) {
            case DELETE_ID:
                Toast.makeText(
                        getApplicationContext(),
                        "Borrando " + current.getNombre(),
                        Toast.LENGTH_LONG).show();
                mParcelaViewModel.delete(current);
                return true;
            case EDIT_ID:
                parelaEdit(current);
                return true;
        }
        return super.onContextItemSelected(item);
    }

    private void createParela() {
        mStartCreateNote.launch(new Intent(this, ParcelaEdit.class));
    }

    ActivityResultLauncher<Intent> mStartCreateNote = newActivityResultLauncher(new ExecuteActivityResult() {
        @Override
        public void process(Bundle extras, Parcela parcela) {
            mParcelaViewModel.insert(parcela);
        }
    });

    ActivityResultLauncher<Intent> newActivityResultLauncher(ExecuteActivityResult executable) {
        return registerForActivityResult(
                new StartActivityForResult(),
                result -> {
                    if (result.getResultCode() == RESULT_OK) {
                        Bundle extras = result.getData().getExtras();
                        Parcela parcela = new Parcela(extras.getString(ParcelaEdit.PARCELA_NOMBRE),
                                                      extras.getString(ParcelaEdit.PARCELA_DESCRIPCION),
                                                      extras.getInt(ParcelaEdit.PARCELA_MAX_OCUPANTES),
                                                      extras.getFloat(ParcelaEdit.PARCELA_PRECIO_P_PERSONA));
                        executable.process(extras, parcela);
                    }
                });
    }

    private void parelaEdit(Parcela current) {
        Intent intent = new Intent(this, ParcelaEdit.class);
        intent.putExtra(ParcelaEdit.PARCELA_NOMBRE, current.getNombre());
        intent.putExtra(ParcelaEdit.PARCELA_DESCRIPCION, current.getDescripcion());
        intent.putExtra(ParcelaEdit.PARCELA_MAX_OCUPANTES, current.getMaxOcupantes());
        intent.putExtra(ParcelaEdit.PARCELA_PRECIO_P_PERSONA, current.getPrecioPorPersona());
        intent.putExtra(ParcelaEdit.PARCELA_ID, current.getIdParcela());
        mStartUpdateParcela.launch(intent);
    }

    ActivityResultLauncher<Intent> mStartUpdateParcela = newActivityResultLauncher(new ExecuteActivityResult() {
        @Override
        public void process(Bundle extras, Parcela parcela) {
            int id = extras.getInt(ParcelaEdit.PARCELA_ID);
            parcela.setIdParcela(id);
            mParcelaViewModel.update(parcela);
        }
    });

}

interface ExecuteActivityResult {
    void process(Bundle extras, Parcela parcela);
}