package es.unizar.eina.M27_camping.ui;

import android.view.ContextMenu;
import android.view.LayoutInflater;
import android.view.Menu;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ImageView;
import android.widget.TextView;

import androidx.recyclerview.widget.RecyclerView;

import es.unizar.eina.M27_camping.R;
import es.unizar.eina.M27_camping.database.Parcela;

class ParcelaViewHolder extends RecyclerView.ViewHolder implements View.OnCreateContextMenuListener {
    private final TextView mParcelaNameView;
    private final TextView mParcelaOcuppantsView;
    private final TextView mParcelaPriceView;

    private final ImageView mEditIcon;
    private final ImageView mDeleteIcon;

    private ParcelaListAdapter.OnEditClickListener editClickListener;
    private ParcelaListAdapter.OnDeleteClickListener deleteClickListener;

    private ParcelaViewHolder(View itemView, ParcelaListAdapter.OnEditClickListener editListener,
                              ParcelaListAdapter.OnDeleteClickListener deleteListener) {
        super(itemView);
        mParcelaNameView = itemView.findViewById(R.id.textView);
        mParcelaOcuppantsView = itemView.findViewById(R.id.parcela_occupants);
        mParcelaPriceView = itemView.findViewById(R.id.parcela_price);
        mEditIcon = itemView.findViewById(R.id.edit_icon);
        mDeleteIcon = itemView.findViewById(R.id.delete_icon);

        this.editClickListener = editListener;
        this.deleteClickListener = deleteListener;

        itemView.setOnCreateContextMenuListener(this);
    }

    public void bind(Parcela parcela) {
        mParcelaNameView.setText(parcela.getNombre());
        mParcelaOcuppantsView.setText("Nº Ocupantes: " + parcela.getMaxOcupantes());
        mParcelaPriceView.setText("Precio/Persona: " + String.format("%.2f€", parcela.getPrecioPorPersona()));
        // Configurar evento de clic para editar
        mEditIcon.setOnClickListener(v -> {
            if (editClickListener != null) {
                editClickListener.onEditClick(parcela);
            }
        });

        // Configurar evento de clic para eliminar
        mDeleteIcon.setOnClickListener(v -> {
            if (deleteClickListener != null) {
                deleteClickListener.onDeleteClick(parcela);
            }
        });
    }



    static ParcelaViewHolder create(ViewGroup parent, ParcelaListAdapter.OnEditClickListener editListener,
                                    ParcelaListAdapter.OnDeleteClickListener deleteListener) {
        View view = LayoutInflater.from(parent.getContext())
                .inflate(R.layout.recyclerview_item, parent, false);
        return new ParcelaViewHolder(view, editListener, deleteListener);
    }


    public void onCreateContextMenu(ContextMenu menu, View v,
                                    ContextMenu.ContextMenuInfo menuInfo) {
        //super.onCreateContextMenu(menu, v, menuInfo);
        menu.add(Menu.NONE, ParcelasListar.DELETE_ID, Menu.NONE, R.string.menu_delete);
        menu.add(Menu.NONE, ParcelasListar.EDIT_ID, Menu.NONE, R.string.menu_edit);
    }


}
