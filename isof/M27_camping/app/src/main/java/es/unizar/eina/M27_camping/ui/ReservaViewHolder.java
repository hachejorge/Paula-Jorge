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
import es.unizar.eina.M27_camping.database.Reserva;

class ReservaViewHolder extends RecyclerView.ViewHolder implements View.OnCreateContextMenuListener {
    private final TextView mReservaNombre;
    private final TextView mReservaTlf;
    private final TextView mReservaFechaEntrada;
    private final TextView mReservaFechaSalida;
    private final TextView mReservaPrecioTotal;

    private final ImageView mEditIcon;
    private final ImageView mDeleteIcon;

    //private ReservaListAdapter.OnEditClickListener editClickListener;
    //private ReservaListAdapter.OnDeleteClickListener deleteClickListener;

    private ReservaViewHolder(View itemView) {
        super(itemView);

        mReservaNombre = itemView.findViewById(R.id.nomCliente);
        mReservaTlf = itemView.findViewById(R.id.tlfCliente);
        mReservaFechaEntrada = itemView.findViewById(R.id.fEntrada);
        mReservaFechaSalida = itemView.findViewById(R.id.fSalida);
        mReservaPrecioTotal = itemView.findViewById(R.id.precioTotal);
        mEditIcon = itemView.findViewById(R.id.edit_icon_reserva);
        mDeleteIcon = itemView.findViewById(R.id.delete_icon_reserva);
        //this.editClickListener = editListener;
        //this.deleteClickListener = deleteListener;

        itemView.setOnCreateContextMenuListener(this);
    }

    public void bind(Reserva reserva) {
        mReservaNombre.setText(reserva.getNomCliente());
        mReservaTlf.setText(reserva.getTlfCliente());
        mReservaFechaEntrada.setText(reserva.getFechaEntrada());
        mReservaFechaSalida.setText(reserva.getFechaSalida());
        //mReservaPrecioTotal.setText(reserva.g);

        /**mEditIcon.setOnClickListener(v -> {
            if (editClickListener != null) {
                //editClickListener.onEditClick(reserva);
            }
        });

        // Configurar evento de clic para eliminar
        mDeleteIcon.setOnClickListener(v -> {
            if (deleteClickListener != null) {
                //deleteClickListener.onDeleteClick(reserva);
            }
        });*/
    }

    static ReservaViewHolder create(ViewGroup parent) {
        View view = LayoutInflater.from(parent.getContext())
                .inflate(R.layout.recyclerview_reserva, parent, false);
        return new ReservaViewHolder(view);
    }


    public void onCreateContextMenu(ContextMenu menu, View v,
                                    ContextMenu.ContextMenuInfo menuInfo) {
        //super.onCreateContextMenu(menu, v, menuInfo);
        menu.add(Menu.NONE, ReservasListar.DELETE_ID, Menu.NONE, R.string.menu_delete);
        menu.add(Menu.NONE, ReservasListar.EDIT_ID, Menu.NONE, R.string.menu_edit);
    }


}
