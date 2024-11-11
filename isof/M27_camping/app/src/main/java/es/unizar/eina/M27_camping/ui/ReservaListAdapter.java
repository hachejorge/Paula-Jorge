package es.unizar.eina.M27_camping.ui;

import android.view.View;
import android.view.ViewGroup;

import androidx.annotation.NonNull;
import androidx.recyclerview.widget.DiffUtil;
import androidx.recyclerview.widget.ListAdapter;

import es.unizar.eina.M27_camping.database.Parcela;
import es.unizar.eina.M27_camping.database.Reserva;

public class ReservaListAdapter extends ListAdapter<Reserva, ReservaViewHolder> {
    private int position;
    private ReservaListAdapter.OnEditClickListener editClickListener;
    private ReservaListAdapter.OnDeleteClickListener deleteClickListener;
    private ReservaListAdapter.OnClickListenerSend sendClickListener;

    public int getPosition() {
        return position;
    }

    public void setPosition(int position) {
        this.position = position;
    }

    public ReservaListAdapter(@NonNull DiffUtil.ItemCallback<Reserva> diffCallback) {
        super(diffCallback);
    }

    // Define interfaces para clics de editar y eliminar
    public interface OnEditClickListener {
        void onEditClick(Reserva reserva);
    }

    public interface OnDeleteClickListener {
        void onDeleteClick(Reserva reserva);
    }

    public interface OnClickListenerSend {
        void onClickSend(Reserva reserva);
    }

    public void setOnEditClickListener(ReservaListAdapter.OnEditClickListener listener) {
        this.editClickListener = listener;
    }

    public void setOnDeleteClickListener(ReservaListAdapter.OnDeleteClickListener listener) {
        this.deleteClickListener = listener;
    }

    public void setSendClickListener(ReservaListAdapter.OnClickListenerSend listener) {
        this.sendClickListener = listener;
    }

    @Override
    public ReservaViewHolder onCreateViewHolder(ViewGroup parent, int viewType) {
        return ReservaViewHolder.create(parent, this.editClickListener, this.deleteClickListener, this.sendClickListener);
    }

    public Reserva getCurrent() {
        return getItem(getPosition());
    }

    @Override
    public void onBindViewHolder(ReservaViewHolder holder, int position) {

        Reserva current = getItem(position);
        holder.bind(current);

        holder.itemView.setOnLongClickListener(new View.OnLongClickListener() {
            @Override
            public boolean onLongClick(View v) {
                setPosition(holder.getAdapterPosition());
                return false;
            }
        });
    }


    static class ReservaDiff extends DiffUtil.ItemCallback<Reserva> {

        @Override
        public boolean areItemsTheSame(@NonNull Reserva oldItem, @NonNull Reserva newItem) {
            //android.util.Log.d ( "ReservaDiff" , "areItemsTheSame " + oldItem.getId() + " vs " + newItem.getId() + " " +  (oldItem.getId() == newItem.getId()));
            return oldItem.getIdReserva() == newItem.getIdReserva();
        }

        @Override
        public boolean areContentsTheSame(@NonNull Reserva oldItem, @NonNull Reserva newItem) {
            //android.util.Log.d ( "ReservaDiff" , "areContentsTheSame " + oldItem.getTitle() + " vs " + newItem.getTitle() + " " + oldItem.getTitle().equals(newItem.getTitle()));
            // We are just worried about differences in visual representation, i.e. changes in the title
            return oldItem.getNomCliente().equals(newItem.getNomCliente());
        }
    }
}
