package es.unizar.eina.M27_camping.ui;

import androidx.appcompat.app.AppCompatActivity;

import android.content.Intent;
import android.os.Bundle;
import android.text.TextUtils;
import android.widget.Button;
import android.widget.EditText;
import android.widget.Toast;

import es.unizar.eina.M27_camping.R;

/** Pantalla utilizada para la creación o edición de una nota */
public class NoteEdit extends AppCompatActivity {

    public static final String NOTE_TITLE = "title";
    public static final String NOTE_BODY = "body";
    public static final String NOTE_ID = "id";

    private EditText mTitleText;

    private EditText mBodyText;

    private Integer mRowId;

    Button mSaveButton;

    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_noteedit);

        mTitleText = findViewById(R.id.title);
        mBodyText = findViewById(R.id.body);

        mSaveButton = findViewById(R.id.button_save);
        mSaveButton.setOnClickListener(view -> {
            Intent replyIntent = new Intent();
            if (TextUtils.isEmpty(mTitleText.getText())) {
                setResult(RESULT_CANCELED, replyIntent);
                Toast.makeText(getApplicationContext(), R.string.empty_not_saved, Toast.LENGTH_LONG).show();
            } else {
                replyIntent.putExtra(NoteEdit.NOTE_TITLE, mTitleText.getText().toString());
                replyIntent.putExtra(NoteEdit.NOTE_BODY, mBodyText.getText().toString());
                if (mRowId!=null) {
                    replyIntent.putExtra(NoteEdit.NOTE_ID, mRowId.intValue());
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
            mTitleText.setText(extras.getString(NoteEdit.NOTE_TITLE));
            mBodyText.setText(extras.getString(NoteEdit.NOTE_BODY));
            mRowId = extras.getInt(NoteEdit.NOTE_ID);
        }
    }

}
