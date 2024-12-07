package es.unizar.eina.M27_camping.ui;

import android.content.Intent;
import android.os.Bundle;
import android.text.TextUtils;
import android.widget.Button;
import android.widget.EditText;
import android.widget.Toast;
import android.app.DatePickerDialog;


import androidx.appcompat.app.AlertDialog;
import androidx.appcompat.app.AppCompatActivity;
import androidx.lifecycle.LiveData;
import androidx.lifecycle.ViewModelProvider;
import androidx.recyclerview.widget.LinearLayoutManager;
import androidx.recyclerview.widget.RecyclerView;

import java.util.Calendar;
import java.util.List;

import es.unizar.eina.M27_camping.R;
import es.unizar.eina.M27_camping.database.ParcelaReservada;

/** Pantalla utilizada para la creación o edición de una reserva */
public class ReservaEdit extends AppCompatActivity {

    public static final String RESERVA_NOMCLIENTE = "nomCliente";
    public static final String RESERVA_TLFCLIENTE = "tlfCliente";
    public static final String RESERVA_FECHAENTRADA = "fechaEntrada";
    public static final String RESERVA_FECHASALIDA = "fechaSalida";
    public static final String RESERVA_ID = "idReserva";

    private EditText mNomClienteText;

    private EditText mTlfClienteText;

    private EditText mFechaEntradaText;

    private EditText mFechaSalidaText;

    private Integer mRowId;

    Button mSaveButton;

    private ParcelaReservadaViewModel mParcelaReservadasViewModel;
    private RecyclerView mRecyclerViewParcelasReservadas;
    private ParcelaReservadaListAdapter mParcelasReservadasAdapter;

    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_reservaedit);

        mNomClienteText = findViewById(R.id.nomCliente);
        mTlfClienteText = findViewById(R.id.tlfCliente);
        mFechaEntradaText = findViewById(R.id.fEntrada);
        mFechaSalidaText = findViewById(R.id.fSalida);

        mRecyclerViewParcelasReservadas = findViewById(R.id.recyclerview_parcelas_reservadas);
        mParcelasReservadasAdapter = new ParcelaReservadaListAdapter(new ParcelaReservadaListAdapter.ParcelaReservadaDiff());
        mRecyclerViewParcelasReservadas.setAdapter(mParcelasReservadasAdapter);
        mRecyclerViewParcelasReservadas.setLayoutManager(new LinearLayoutManager(this));

        mParcelaReservadasViewModel = new ViewModelProvider(this).get(ParcelaReservadaViewModel.class);

        // Añadir el selector de fecha para mFechaEntradaText
        mFechaEntradaText.setOnClickListener(v -> mostrarDatePicker(mFechaEntradaText));

        // Añadir el selector de fecha para mFechaSalidaText
        mFechaSalidaText.setOnClickListener(v -> mostrarDatePicker(mFechaSalidaText));


        int id_reserva = getIntent().getIntExtra(ReservaEdit.RESERVA_ID, -1);

        if(id_reserva != -1) {
            mParcelaReservadasViewModel.getAllParcelasPorReserva(id_reserva).observe(this, parcelaReservadas -> {
                if (parcelaReservadas != null ) {
                    mParcelasReservadasAdapter.submitList(parcelaReservadas);
                }
            });
        }

        mParcelasReservadasAdapter.setOnAumentarClickListener(parcelaReservada -> {
            // Hay que chequear el máximo
        });

        mParcelasReservadasAdapter.setOnDisminuirClickListener(parcelaReservada -> {
            int ocupantesActuales = parcelaReservada.getNumOcupantes();
            int position = getIndexOfParcela(parcelaReservada);
            // Tiene dos o más ocupantes
            if(ocupantesActuales > 1){
                parcelaReservada.setNumOcupantes(ocupantesActuales - 1);
                mParcelaReservadasViewModel.update(parcelaReservada);
                mParcelasReservadasAdapter.updateItem(position, parcelaReservada);
                Toast.makeText(this, "Ocupantes disminuidos a " + parcelaReservada.getNumOcupantes(), Toast.LENGTH_SHORT).show();
            }
            // No puede tener menos de 1 ocupante
            else{
                Toast.makeText(this, "No se puede disminuir más, el mínimo es de un ocupante", Toast.LENGTH_SHORT).show();
            }
        });

        mParcelasReservadasAdapter.setOnDeleteClickListener(parcelaReservada -> {
            new AlertDialog.Builder(this)
                    .setTitle("Confirmar eliminación")
                    .setMessage("¿Estás seguro de que quieres eliminar " + parcelaReservada.getNomParcela() + " de la reserva #" + parcelaReservada.getIdReservaPR() + "?")
                    .setPositiveButton("Sí", (dialog, which) -> {
                        // Eliminar la parcela si el usuario confirma
                        Toast.makeText(
                                getApplicationContext(),
                                "Borrando " + parcelaReservada.getNomParcela(),
                                Toast.LENGTH_LONG).show();
                        mParcelaReservadasViewModel.delete(parcelaReservada);
                    })
                    .setNegativeButton("No", (dialog, which) -> {
                        // Cancelar la eliminación
                        dialog.dismiss();
                    })
                    .show();
        });

        mSaveButton = findViewById(R.id.button_save);
        mSaveButton.setOnClickListener(view -> {
            Intent replyIntent = new Intent();
            if (TextUtils.isEmpty(mNomClienteText.getText())) {
                setResult(RESULT_CANCELED, replyIntent);
                Toast.makeText(getApplicationContext(), R.string.reserva_not_saved, Toast.LENGTH_LONG).show();
            }// Si la fecha de entrada es igual o posterior a la fecha de salida
            else if (mFechaEntradaText.getText().toString().compareTo(mFechaSalidaText.getText().toString()) >= 0) {
                setResult(RESULT_CANCELED, replyIntent);
                Toast.makeText(getApplicationContext(), R.string.fecha_entrada_salida_error, Toast.LENGTH_LONG).show();
            } else {

                replyIntent.putExtra(ReservaEdit.RESERVA_NOMCLIENTE, mNomClienteText.getText().toString());
                replyIntent.putExtra(ReservaEdit.RESERVA_FECHAENTRADA, mFechaEntradaText.getText().toString());
                replyIntent.putExtra(ReservaEdit.RESERVA_FECHASALIDA, mFechaSalidaText.getText().toString());

                // Convertir Telefono a un entero
                try {
                    int numTlf = Integer.parseInt(mTlfClienteText.getText().toString());
                    replyIntent.putExtra(ReservaEdit.RESERVA_TLFCLIENTE, numTlf);
                } catch (NumberFormatException e) {
                    Toast.makeText(getApplicationContext(), "Número de teléfono inválido", Toast.LENGTH_LONG).show();
                    return;
                }

                if (mRowId!=null) {
                    replyIntent.putExtra(ReservaEdit.RESERVA_ID, mRowId.intValue());
                }
                setResult(RESULT_OK, replyIntent);
            }
            finish();
        });

        populateFields();

    }

    private void populateFields () {
        mRowId = null;
        Bundle extras = getIntent().getExtras();
        if (extras!=null) {
            mNomClienteText.setText(extras.getString(ReservaEdit.RESERVA_NOMCLIENTE));

            int tlf = extras.getInt(ReservaEdit.RESERVA_TLFCLIENTE, 0);
            mTlfClienteText.setText(String.valueOf(tlf));

            mFechaEntradaText.setText(extras.getString(ReservaEdit.RESERVA_FECHAENTRADA));
            mFechaSalidaText.setText(extras.getString(ReservaEdit.RESERVA_FECHASALIDA));
            mRowId = extras.getInt(ReservaEdit.RESERVA_ID);
        }
    }

    // Método para mostrar el DatePickerDialog
    private void mostrarDatePicker(EditText editText) {
        final Calendar calendar = Calendar.getInstance();
        int year = calendar.get(Calendar.YEAR);
        int month = calendar.get(Calendar.MONTH);
        int day = calendar.get(Calendar.DAY_OF_MONTH);

        DatePickerDialog datePickerDialog = new DatePickerDialog(this, (view, year1, month1, dayOfMonth) -> {
            // Formatear la fecha como AAAA-MM-DD
            String fechaSeleccionada = String.format("%04d-%02d-%02d", year1, month1 + 1, dayOfMonth);
            editText.setText(fechaSeleccionada);
        }, year, month, day);

        datePickerDialog.show();
    }

    private int getIndexOfParcela(ParcelaReservada parcelaReservada) {
        List<ParcelaReservada> currentList = mParcelasReservadasAdapter.getCurrentList();
        for (int i = 0; i < currentList.size(); i++) {
            if (currentList.get(i).getIdParcelaPR() == parcelaReservada.getIdParcelaPR()) {
                return i; // Devolver la posición si coincide el ID
            }
        }
        return -1; // Si no se encuentra
    }



}
