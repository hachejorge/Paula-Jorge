package es.unizar.eina.M27_camping.ui;

import android.content.Intent;
import android.os.Bundle;
import android.text.TextUtils;
import android.widget.Button;
import android.widget.EditText;
import android.widget.Toast;
import android.app.DatePickerDialog;


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

        mParcelaReservadasViewModel.getAllParcelasReservadas().observe(this, parcelaReservadas -> {
            if (parcelaReservadas != null ) {
                mParcelasReservadasAdapter.submitList(parcelaReservadas);
            }
        });

        mSaveButton = findViewById(R.id.button_save);
        mSaveButton.setOnClickListener(view -> {
            Intent replyIntent = new Intent();
            if (TextUtils.isEmpty(mNomClienteText.getText())) {
                setResult(RESULT_CANCELED, replyIntent);
                Toast.makeText(getApplicationContext(), R.string.reserva_not_saved, Toast.LENGTH_LONG).show();
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

}
