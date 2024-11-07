package es.unizar.eina.M27_camping.ui;

import androidx.appcompat.app.AppCompatActivity;

import android.content.Intent;
import android.os.Bundle;
import android.text.TextUtils;
import android.widget.Button;
import android.widget.EditText;
import android.widget.Toast;

import es.unizar.eina.M27_camping.R;

/** Pantalla utilizada para la creación o edición de una parcela */
public class ParcelaEdit extends AppCompatActivity {

    public static final String PARCELA_NOMBRE = "nombre";
    public static final String PARCELA_DESCRIPCION = "descripcion";
    public static final String PARCELA_MAX_OCUPANTES = "maxOcupantes";
    public static final String PARCELA_PRECIO_P_PERSONA = "precioPorPersona";
    public static final String PARCELA_ID = "idParcela";

    private EditText mNombreText;

    private EditText mDescripcionText;

    private EditText mMaxOcupantesText;

    private EditText mPrecioPorPersonaText;

    private Integer mRowId;

    Button mSaveButton;

    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_parcelaedit);

        mNombreText = findViewById(R.id.nombre);
        mDescripcionText = findViewById(R.id.descripcion);
        mMaxOcupantesText = findViewById(R.id.maxOcupantes);
        mPrecioPorPersonaText = findViewById(R.id.precioPorPersona);

        mSaveButton = findViewById(R.id.button_save);
        mSaveButton.setOnClickListener(view -> {
            Intent replyIntent = new Intent();
            if (TextUtils.isEmpty(mNombreText.getText())) {
                setResult(RESULT_CANCELED, replyIntent);
                Toast.makeText(getApplicationContext(), R.string.empty_not_saved, Toast.LENGTH_LONG).show();
            } else {
                replyIntent.putExtra(ParcelaEdit.PARCELA_NOMBRE, mNombreText.getText().toString());
                replyIntent.putExtra(ParcelaEdit.PARCELA_DESCRIPCION, mDescripcionText.getText().toString());
                replyIntent.putExtra(ParcelaEdit.PARCELA_MAX_OCUPANTES, mMaxOcupantesText.getText().toString());
                replyIntent.putExtra(ParcelaEdit.PARCELA_PRECIO_P_PERSONA, mPrecioPorPersonaText.getText().toString());
                if (mRowId!=null) {
                    replyIntent.putExtra(ParcelaEdit.PARCELA_ID, mRowId.intValue());
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
            mNombreText.setText(extras.getString(ParcelaEdit.PARCELA_NOMBRE));
            mDescripcionText.setText(extras.getString(ParcelaEdit.PARCELA_DESCRIPCION));
            mMaxOcupantesText.setText(extras.getString(ParcelaEdit.PARCELA_MAX_OCUPANTES));
            mPrecioPorPersonaText.setText(extras.getString(ParcelaEdit.PARCELA_PRECIO_P_PERSONA));
            mRowId = extras.getInt(ParcelaEdit.PARCELA_ID);
        }
    }

}
