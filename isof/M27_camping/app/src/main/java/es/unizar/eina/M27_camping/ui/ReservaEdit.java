package es.unizar.eina.M27_camping.ui;

import android.content.Intent;
import android.os.Bundle;
import android.text.TextUtils;
import android.view.View;
import android.widget.AdapterView;
import android.widget.ArrayAdapter;
import android.widget.Button;
import android.widget.EditText;
import android.widget.Spinner;
import android.widget.Toast;
import android.app.DatePickerDialog;


import androidx.appcompat.app.AlertDialog;
import androidx.appcompat.app.AppCompatActivity;
import androidx.lifecycle.LiveData;
import androidx.lifecycle.ViewModelProvider;
import androidx.recyclerview.widget.LinearLayoutManager;
import androidx.recyclerview.widget.RecyclerView;

import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.time.LocalDate;
import java.time.format.DateTimeFormatter;
import java.time.temporal.ChronoUnit;
import java.util.ArrayList;
import java.util.Calendar;
import java.util.Collections;
import java.util.Date;
import java.util.List;
import java.util.concurrent.Executors;
import java.util.concurrent.atomic.AtomicReference;

import es.unizar.eina.M27_camping.R;
import es.unizar.eina.M27_camping.database.Parcela;
import es.unizar.eina.M27_camping.database.ParcelaReservada;

/** Pantalla utilizada para la creación o edición de una reserva */
public class ReservaEdit extends AppCompatActivity {

    public static final String RESERVA_NOMCLIENTE = "nomCliente";
    public static final String RESERVA_TLFCLIENTE = "tlfCliente";
    public static final String RESERVA_FECHAENTRADA = "fechaEntrada";
    public static final String RESERVA_FECHASALIDA = "fechaSalida";
    public static final String RESERVA_ID = "idReserva";
    public static final String RESERVA_PRECIO = "precioTotal";

    private EditText mNomClienteText;

    private EditText mTlfClienteText;

    private EditText mFechaEntradaText;

    private EditText mFechaSalidaText;

    private Integer mRowId;

    Button mSaveButton;
    Button mAddButton;

    private ParcelaReservadaViewModel mParcelaReservadasViewModel;
    private ParcelaViewModel mParcelaViewModel;
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
        mParcelaViewModel = new ViewModelProvider(this).get(ParcelaViewModel.class);

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

        //SPINNER PARCELA RESERVADA
        Spinner spinnerReserva = findViewById(R.id.parcelas_disponibles);

        // Observar las parcelas disponibles y las reservadas
        mParcelaViewModel.getParcelasPorOcupantes().observe(this, parcelasDisponibles -> {
            mParcelaReservadasViewModel.getAllParcelasPorReserva(id_reserva).observe(this, parcelasReservadas -> {
                if (parcelasDisponibles != null && !parcelasDisponibles.isEmpty()) {
                    // Filtrar las parcelas que no están reservadas
                    List<String> nombresParcelasFiltradas = new ArrayList<>();
                    List<Integer> idsParcelasReservadas = new ArrayList<>();

                    // Obtener los IDs de las parcelas reservadas
                    for (ParcelaReservada parcelaReservada : parcelasReservadas) {
                        idsParcelasReservadas.add(parcelaReservada.getIdParcelaPR());
                    }

                    // Filtrar las parcelas disponibles
                    for (Parcela parcela : parcelasDisponibles) {
                        if (!idsParcelasReservadas.contains(parcela.getIdParcela())) {
                            nombresParcelasFiltradas.add(parcela.getNombre());
                        }
                    }

                    // Configurar el adaptador del Spinner con las parcelas filtradas
                    ArrayAdapter<String> adapter = new ArrayAdapter<>(this, android.R.layout.simple_spinner_item, nombresParcelasFiltradas);
                    adapter.setDropDownViewResource(android.R.layout.simple_spinner_dropdown_item);
                    spinnerReserva.setAdapter(adapter);
                } else {
                    Toast.makeText(this, "No hay parcelas disponibles", Toast.LENGTH_SHORT).show();
                }
            });
        });


        spinnerReserva.setOnItemSelectedListener(new AdapterView.OnItemSelectedListener() {
            @Override
            public void onItemSelected(AdapterView<?> parent, View view, int position, long id) {
                // Cuando se haga click en el boton "+" se añadirá la parcela seleccionada
            }
            @Override
            public void onNothingSelected(AdapterView<?> parent) {
                // No hacer nada si no se selecciona nada
            }
        });


        // BOTON ADD RESERVA
        mAddButton = findViewById(R.id.button_addParcela);
        mAddButton.setOnClickListener(view -> {
            Intent replyIntent = new Intent();

            String nomParcela = (String) spinnerReserva.getSelectedItem();
            mParcelaViewModel.getParcelaPorNombre(nomParcela).observe(this, parcela -> {
                ParcelaReservada pr = new ParcelaReservada(id_reserva,
                        parcela.getIdParcela(), nomParcela, 1);
                mParcelaReservadasViewModel.insert(pr);

                mParcelaReservadasViewModel.getAllParcelasPorReserva(id_reserva).observe(this, parcelaReservadas -> {
                    if (parcelaReservadas != null ) {
                        mParcelasReservadasAdapter.submitList(parcelaReservadas);
                    }
                });
            });

        });

        mParcelasReservadasAdapter.setOnAumentarClickListener(parcelaReservada -> {
            // Hay que chequear el máximo
            int ocupantesActuales = parcelaReservada.getNumOcupantes();
            int position = getIndexOfParcela(parcelaReservada);
            mParcelaViewModel.getParcelaPorNombre(parcelaReservada.getNomParcela()).observe(this, parcela -> {
                //Tiene menos ocupantes del máximo
                if (ocupantesActuales < parcela.getMaxOcupantes()) {
                    parcelaReservada.setNumOcupantes(ocupantesActuales + 1);
                    mParcelaReservadasViewModel.update(parcelaReservada);
                    mParcelasReservadasAdapter.updateItem(position, parcelaReservada);
                    mParcelaReservadasViewModel.getAllParcelasPorReserva(id_reserva).observe(this, parcelasReservadas -> {
                        if (parcelasReservadas != null ) {
                            mParcelasReservadasAdapter.submitList(parcelasReservadas);
                        }
                    });
                    mParcelaReservadasViewModel.getAllParcelasPorReserva(id_reserva).removeObservers(this);
                    Toast.makeText(this, "Ocupantes aumentados a " + parcelaReservada.getNumOcupantes(), Toast.LENGTH_SHORT).show();
                }
                //Hay el mismo número de ocupantes que el máximo
                else {
                    Toast.makeText(this, "Ocupantes máximos alcanzados", Toast.LENGTH_SHORT).show();
                }
            });

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

            // Validación del nombre del cliente
            if (TextUtils.isEmpty(mNomClienteText.getText())) {
                setResult(RESULT_CANCELED, replyIntent);
                Toast.makeText(getApplicationContext(), R.string.reserva_not_saved, Toast.LENGTH_LONG).show();
                return;
            }

            // Validación de fechas
            if (mFechaEntradaText.getText().toString().compareTo(mFechaSalidaText.getText().toString()) >= 0) {
                setResult(RESULT_CANCELED, replyIntent);
                Toast.makeText(getApplicationContext(), R.string.fecha_entrada_salida_error, Toast.LENGTH_LONG).show();
                return;
            }



                // Manejo de precio total en un hilo de fondo
            Executors.newSingleThreadExecutor().execute(() -> {
                float totalPrice = 0.0f;
                boolean haySolape = false;
                String fEntrada = mFechaEntradaText.getText().toString();
                String fSalida = mFechaSalidaText.getText().toString();
                long diasDeEstancia = calcularDiasEntreFechas(fEntrada, fSalida);
                // Obtiene las parcelas reservadas
                List<ParcelaReservada> parcelaReservadas = mParcelaReservadasViewModel.getListaParcelasPorReserva(id_reserva);
                List<ParcelaReservada> parcelasSolapadas = new ArrayList<>();
                for (ParcelaReservada pr : parcelaReservadas) {
                    Parcela parcela = mParcelaViewModel.getParcelaPorId(pr.getIdParcelaPR());
                    if (parcela != null) {
                        totalPrice += parcela.getPrecioPorPersona() * pr.getNumOcupantes() * diasDeEstancia;
                    }
                    // Se produce solape
                    if(mParcelaReservadasViewModel.getReservasConSolape(id_reserva, pr.getIdParcelaPR(), fEntrada, fSalida) > 0){
                        haySolape = true;
                        parcelasSolapadas.add(pr);
                    }
                }

                // Si se produce solape elimina las parcelas de la reserva e informa sobre el solap
                if(haySolape){
                    for(ParcelaReservada pr : parcelasSolapadas){
                        mParcelaReservadasViewModel.delete(pr);
                    }

                    // Regresar al hilo principal para mostrar un mensaje al usuario
                    runOnUiThread(() -> {
                        Toast.makeText(getApplicationContext(), "No se puede guardar la reserva debido a solapamientos", Toast.LENGTH_LONG).show();
                        setResult(RESULT_CANCELED, replyIntent);
                        finish(); // Finalizar la actividad sin guardar
                    });
                    return;
                }

                // Regresa al hilo principal para manejar los datos y finalizar la actividad
                float finalTotalPrice = totalPrice;
                runOnUiThread(() -> {
                    replyIntent.putExtra(ReservaEdit.RESERVA_PRECIO, finalTotalPrice);
                    replyIntent.putExtra(ReservaEdit.RESERVA_NOMCLIENTE, mNomClienteText.getText().toString());
                    replyIntent.putExtra(ReservaEdit.RESERVA_FECHAENTRADA, fEntrada);
                    replyIntent.putExtra(ReservaEdit.RESERVA_FECHASALIDA, fSalida);

                    // Validación del número de teléfono
                    try {
                        int numTlf = Integer.parseInt(mTlfClienteText.getText().toString());
                        replyIntent.putExtra(ReservaEdit.RESERVA_TLFCLIENTE, numTlf);
                    } catch (NumberFormatException e) {
                        Toast.makeText(getApplicationContext(), "Número de teléfono inválido", Toast.LENGTH_LONG).show();
                        return;
                    }

                    if (mRowId != null) {
                        replyIntent.putExtra(ReservaEdit.RESERVA_ID, mRowId.intValue());
                    }

                    // Ejecutar la actualización en segundo plano
                    Executors.newSingleThreadExecutor().execute(() -> {
                        mParcelaReservadasViewModel.actualizarReservaPendiente();
                    });

                    setResult(RESULT_OK, replyIntent);
                    finish();
                });
            });
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

    private static long calcularDiasEntreFechas(String fechaInicio, String fechaFin) {
        // Formato de las fechas
        SimpleDateFormat formatter = new SimpleDateFormat("yyyy-MM-dd");

        try {
            // Parsear las cadenas de entrada a objetos Date
            Date inicio = formatter.parse(fechaInicio);
            Date fin = formatter.parse(fechaFin);

            // Calcular la diferencia en milisegundos
            long diferenciaMilisegundos = fin.getTime() - inicio.getTime();

            // Convertir la diferencia a días
            return diferenciaMilisegundos / (1000 * 60 * 60 * 24);
        } catch (ParseException e) {
            e.printStackTrace();
            return -1; // Devuelve un valor negativo si hay un error
        }
    }
}
