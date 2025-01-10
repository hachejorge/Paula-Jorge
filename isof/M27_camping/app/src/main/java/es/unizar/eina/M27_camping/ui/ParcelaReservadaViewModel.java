package es.unizar.eina.M27_camping.ui;

import android.app.Application;
import android.util.Log;

import androidx.lifecycle.AndroidViewModel;
import androidx.lifecycle.LiveData;

import java.util.List;
import java.util.concurrent.Executors;

import es.unizar.eina.M27_camping.database.ParcelaReservada;
import es.unizar.eina.M27_camping.database.ParcelaReservadaRepository;

public class ParcelaReservadaViewModel extends AndroidViewModel {

    private ParcelaReservadaRepository mRepository;

    public ParcelaReservadaViewModel(Application application) {
        super(application);
        mRepository = new ParcelaReservadaRepository(application);
    }

    LiveData<List<ParcelaReservada>> getAllParcelasPorReserva(int id_reserva) {
        return mRepository.getParcelasFromReserva(id_reserva);
    }

    List<ParcelaReservada> getListaParcelasPorReserva(int id_reserva) {
        return mRepository.getListaParcelasFromReserva(id_reserva);
    }

    void actualizarReservaPendiente(){
        mRepository.actualizarReservasPendientes();
    }

    public void insert(ParcelaReservada parcelaReservada) { mRepository.insert(parcelaReservada); }
    public void update(ParcelaReservada parcelaReservada) { mRepository.update(parcelaReservada); }
    public void delete(ParcelaReservada parcelaReservada) { mRepository.delete(parcelaReservada); }
    public void deleteByParcelaId(int idParcela) {
        Executors.newSingleThreadExecutor().execute(() -> {
            mRepository.deleteByParcelaId(idParcela);
        });
    }

}
