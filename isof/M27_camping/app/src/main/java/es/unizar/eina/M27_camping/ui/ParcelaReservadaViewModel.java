package es.unizar.eina.M27_camping.ui;

import android.app.Application;

import androidx.lifecycle.AndroidViewModel;
import androidx.lifecycle.LiveData;

import java.util.List;

import es.unizar.eina.M27_camping.database.ParcelaReservada;
import es.unizar.eina.M27_camping.database.ParcelaReservadaRepository;

public class ParcelaReservadaViewModel extends AndroidViewModel {

    private ParcelaReservadaRepository mRepository;

    private final LiveData<List<ParcelaReservada>> mAllParcelasReserva;

    public ParcelaReservadaViewModel(Application application) {
        super(application);
        mRepository = new ParcelaReservadaRepository(application);
        mAllParcelasReserva = mRepository.getAllParcelasReservadas();
    }

    LiveData<List<ParcelaReservada>> getAllParcelasReservadas() { return mAllParcelasReserva; }


    public void insert(ParcelaReservada parcelaReservada) { mRepository.insert(parcelaReservada); }
    public void update(ParcelaReservada parcelaReservada) { mRepository.update(parcelaReservada); }
    public void delete(ParcelaReservada parcelaReservada) { mRepository.delete(parcelaReservada); }
}
