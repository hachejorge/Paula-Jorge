package es.unizar.eina.M27_camping.ui;

import android.app.Application;

import androidx.lifecycle.AndroidViewModel;
import androidx.lifecycle.LiveData;

import java.util.List;

import es.unizar.eina.M27_camping.database.Parcela;
import es.unizar.eina.M27_camping.database.ParcelaRepository;

public class ParcelaViewModel extends AndroidViewModel {

    private ParcelaRepository mRepository;

    private final LiveData<List<Parcela>> mAllParcelas;

    public ParcelaViewModel(Application application) {
        super(application);
        mRepository = new ParcelaRepository(application);
        mAllParcelas = mRepository.getAllParcelasPorNombre();
    }

    public LiveData<List<Parcela>> getParcelasPorNombre() {
        return mRepository.getAllParcelasPorNombre();
    }

    public LiveData<List<Parcela>> getParcelasPorPrecio() {
        return mRepository.getAllParcelasPorOcupantes();
    }

    public LiveData<List<Parcela>> getParcelasPorOcupantes() {
        return mRepository.getAllParcelasPorPrecio();
    }


    public void insert(Parcela parcela) { mRepository.insert(parcela); }

    public void update(Parcela parcela) { mRepository.update(parcela); }
    public void delete(Parcela parcela) { mRepository.delete(parcela); }
}
