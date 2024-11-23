package es.unizar.eina.M27_camping.ui;

import android.app.Application;

import androidx.lifecycle.AndroidViewModel;
import androidx.lifecycle.LiveData;

import java.util.List;

import es.unizar.eina.M27_camping.database.ParcelaReservada;
import es.unizar.eina.M27_camping.database.ParcelaReservadaRepository;

public class ParcelaReservadaViewModel extends AndroidViewModel {

    private ParcelaReservadaRepository mRepository;

    //private final LiveData<List<Reserva>> mAllReservas;

    public ParcelaReservadaViewModel(Application application) {
        super(application);
        mRepository = new ParcelaReservadaRepository(application);
        //mAllReservas = mRepository.getAllReservas();
    }

    LiveData<List<ParcelaReservada>> getAllParcelasReservadas() { return mRepository.getAllParcelasReservadas(); }


    public void insert(ParcelaReservada parcelaReservada) { mRepository.insert(parcelaReservada); }
    public void update(ParcelaReservada parcelaReservada) { mRepository.update(parcelaReservada); }
    public void delete(ParcelaReservada parcelaReservada) { mRepository.delete(parcelaReservada); }
}
