package es.unizar.eina.M27_camping.ui;

import android.view.View;
import android.view.ViewGroup;

import androidx.annotation.NonNull;
import androidx.recyclerview.widget.DiffUtil;
import androidx.recyclerview.widget.ListAdapter;

import java.util.List;

import es.unizar.eina.M27_camping.database.ParcelaReservada;

public class ParcelaReservadaListAdapter extends ListAdapter<ParcelaReservada, ParcelaReservadaViewHolder> {

    private int position;
    private List<ParcelaReservada> parcelas;

    private OnAumentarClickListener aumentarClickListener;
    private OnDisminuirClickListener disminuirClickListener;
    private OnDeleteClickListener deleteClickListener;

    // Define interfaces para clics de eliminar, además de aumentar y disminuir los ocupantes
    public interface OnAumentarClickListener {
        void onAumentarClick(ParcelaReservada parcelaReservada);
    }

    public interface OnDisminuirClickListener {
        void onDisminuirClick(ParcelaReservada parcelaReservada);
    }

    public interface OnDeleteClickListener {
        void onDeleteClick(ParcelaReservada parcelaReservada);
    }




    public void setOnAumentarClickListener(OnAumentarClickListener listener) {
        this.aumentarClickListener = listener;
    }

    public void setOnDisminuirClickListener(OnDisminuirClickListener listener) {
        this.disminuirClickListener = listener;
    }

    public void setOnDeleteClickListener(OnDeleteClickListener listener) {
        this.deleteClickListener = listener;
    }



    public int getPosition() {
        return position;
    }

    public void setPosition(int position) {
      this.position = position;
    }

    public ParcelaReservadaListAdapter(@NonNull DiffUtil.ItemCallback<ParcelaReservada> diffCallback) {
        super(diffCallback);
    }

    @Override
    public ParcelaReservadaViewHolder onCreateViewHolder(ViewGroup parent, int viewType) {
        return ParcelaReservadaViewHolder.create(parent, aumentarClickListener, disminuirClickListener, deleteClickListener);
    }

    public ParcelaReservada getCurrent() {
        return getItem(getPosition());
    }

    @Override
    public void onBindViewHolder(ParcelaReservadaViewHolder holder, int position) {

        ParcelaReservada current = getItem(position);
        holder.bind(current);

        holder.itemView.setOnLongClickListener(new View.OnLongClickListener() {
            @Override
            public boolean onLongClick(View v) {
                setPosition(holder.getAdapterPosition());
                return false;
            }
        });
    }


    static class ParcelaReservadaDiff extends DiffUtil.ItemCallback<ParcelaReservada> {

        @Override
        public boolean areItemsTheSame(@NonNull ParcelaReservada oldItem, @NonNull ParcelaReservada newItem) {
            //android.util.Log.d ( "ParcelaDiff" , "areItemsTheSame " + oldItem.getId() + " vs " + newItem.getId() + " " +  (oldItem.getId() == newItem.getId()));
            return oldItem.getIdParcelaPR() == newItem.getIdParcelaPR();
        }

        @Override
        public boolean areContentsTheSame(@NonNull ParcelaReservada oldItem, @NonNull ParcelaReservada newItem) {
            //android.util.Log.d ( "ParcelaDiff" , "areContentsTheSame " + oldItem.getTitle() + " vs " + newItem.getTitle() + " " + oldItem.getTitle().equals(newItem.getTitle()));
            // We are just worried about differences in visual representation, i.e. changes in the title
            return oldItem.getNomParcela().equals(newItem.getNomParcela());
        }
    }
}
