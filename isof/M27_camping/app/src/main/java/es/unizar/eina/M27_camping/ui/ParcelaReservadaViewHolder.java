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
import es.unizar.eina.M27_camping.database.ParcelaReservada;

class ParcelaReservadaViewHolder extends RecyclerView.ViewHolder implements View.OnCreateContextMenuListener {
    private final TextView mParcelaNameView;
    private final TextView mParcelaOcuppantsView;

    private final ImageView mAumentarIcon;
    private final ImageView mDisminuirIcon;
    private final ImageView mDeleteIcon;

    private ParcelaReservadaListAdapter.OnAumentarClickListener aumentarClickListener;
    private ParcelaReservadaListAdapter.OnDisminuirClickListener disminuirClickListener;
    private ParcelaReservadaListAdapter.OnDeleteClickListener deleteClickListener;

    private ParcelaReservadaViewHolder(View itemView, ParcelaReservadaListAdapter.OnAumentarClickListener aumentarListener,
                                       ParcelaReservadaListAdapter.OnDisminuirClickListener disminuirListener,
                                       ParcelaReservadaListAdapter.OnDeleteClickListener deleteListener) {
        super(itemView);
        mParcelaNameView = itemView.findViewById(R.id.textView);
        mParcelaOcuppantsView = itemView.findViewById(R.id.parcela_occupants);
        mAumentarIcon = itemView.findViewById(R.id.edit_icon);
        mDisminuirIcon = itemView.findViewById(R.id.edit_icon);
        mDeleteIcon = itemView.findViewById(R.id.delete_icon);

        this.aumentarClickListener = aumentarListener;
        this.disminuirClickListener = disminuirListener;
        this.deleteClickListener = deleteListener;

        itemView.setOnCreateContextMenuListener(this);
    }

    public void bind(ParcelaReservada parcelaReservada) {
        mParcelaNameView.setText(parcelaReservada.getNomParcela());
        mParcelaOcuppantsView.setText(parcelaReservada.getNumOcupantes());
        // Configurar evento de clic para aumentar
        mAumentarIcon.setOnClickListener(v -> {
            if (aumentarClickListener != null) {
                aumentarClickListener.onAumentarClick(parcelaReservada);
            }
        });

        // Configurar evento de clic para disminuir
        mDisminuirIcon.setOnClickListener(v -> {
            if (disminuirClickListener != null) {
                disminuirClickListener.onDisminuirClick(parcelaReservada);
            }
        });

        // Configurar evento de clic para eliminar
        mDeleteIcon.setOnClickListener(v -> {
            if (deleteClickListener != null) {
                deleteClickListener.onDeleteClick(parcelaReservada);
            }
        });
    }



    static ParcelaReservadaViewHolder create(ViewGroup parent, ParcelaReservadaListAdapter.OnAumentarClickListener aumentarListener,
                                             ParcelaReservadaListAdapter.OnDisminuirClickListener disminuirListener,
                                             ParcelaReservadaListAdapter.OnDeleteClickListener deleteListener) {
        View view = LayoutInflater.from(parent.getContext())
                .inflate(R.layout.recyclerview_item, parent, false);
        return new ParcelaReservadaViewHolder(view, aumentarListener, disminuirListener, deleteListener);
    }

//?????????????????????????
    public void onCreateContextMenu(ContextMenu menu, View v,
                                    ContextMenu.ContextMenuInfo menuInfo) {
        //super.onCreateContextMenu(menu, v, menuInfo);
        menu.add(Menu.NONE, ParcelasListar.DELETE_ID, Menu.NONE, R.string.menu_delete);
        menu.add(Menu.NONE, ParcelasListar.EDIT_ID, Menu.NONE, R.string.menu_edit);
    }



}
