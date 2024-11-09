package es.unizar.eina.M27_camping.ui;

import android.app.Application;

import androidx.lifecycle.AndroidViewModel;
import androidx.lifecycle.LiveData;

import java.util.List;

import es.unizar.eina.M27_camping.database.Reserva;
import es.unizar.eina.M27_camping.database.ReservaRepository;

public class ReservaViewModel extends AndroidViewModel {

    private ReservaRepository mRepository;

    private final LiveData<List<Reserva>> mAllReservas;

    public ReservaViewModel(Application application) {
        super(application);
        mRepository = new ReservaRepository(application);
        mAllReservas = mRepository.getAllReservas();
    }

    LiveData<List<Reserva>> getAllReservas() { return mAllReservas; }

    public void insert(Reserva reserva) { mRepository.insert(reserva); }

    public void update(Reserva reserva) { mRepository.update(reserva); }
    public void delete(Reserva reserva) { mRepository.delete(reserva); }
}
