import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from 'src/environments/environment';
import { Dataset, Ng2Dataset, Sensordata } from '../models/dataset';
import { Device } from '../models/device';

@Injectable({
  providedIn: 'root',
})
export class DatasetsService {
  constructor(private http: HttpClient) {}

  public LoadDatasets() {
    return this.http.get<Dataset[]>(`${environment.apiUrl}/api/datasets`);
  }
  public LoadDataset(reference: string) {
    return this.http.get<Dataset>(`${environment.apiUrl}/api/datasets/${reference}`);
  }
  public LoadDatasetByReference(reference: string) {
    return this.http.get<Sensordata[]>(`${environment.apiUrl}/api/sensordata/${reference}`);
  }
  public LoadCSVDatasetByReference(reference: string) {
    return this.http.get<any[]>(`${environment.apiUrl}/api/exportdataset/${reference}`);
  }
  public LoadAreaChartDatasetByReference(reference: string, limit: number) {
    return this.http.get<any[]>(`${environment.apiUrl}/api/chart/area/${reference}/${limit}`);
  }
  public LoadLineChartDatasetByReference(reference: string, limit: number) {
    return this.http.get<Ng2Dataset>(`${environment.apiUrl}/api/chart/line/${reference}/${limit}`);
  }
  public newDataset(device: Device) {
    return this.http.post<Device>(`${environment.apiUrl}/api/datasets`, device);
  }
  public updateDataset(dataset: Dataset) {
    return this.http.put<Dataset>(`${environment.apiUrl}/api/datasets`, dataset);
  }
  public DeleteDatasetByID(id: number, title: string) {
    return this.http.post<any>(`${environment.apiUrl}/api/datasets/delete`, {
      id,
      title,
    });
  }
}
