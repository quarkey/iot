import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { environment } from "src/environments/environment";
import { Dataset, Sensordata } from "../models/dataset";
import { Device } from "../models/device";

@Injectable({
  providedIn: "root",
})
export class DatasetsService {
  constructor(private http: HttpClient) {}

  public LoadDatasets() {
    return this.http.get<Dataset[]>(`${environment.apiUrl}/api/datasets`);
  }
  public LoadDataset(reference: string) {
    return this.http.get<Dataset>(
      `${environment.apiUrl}/api/datasets/${reference}`
    );
  }
  public LoadDatasetByReference(reference: string) {
    return this.http.get<Sensordata[]>(
      `${environment.apiUrl}/api/sensordata/${reference}`
    );
  }
  public LoadAreaChartDatasetByReference(reference: string) {
    return this.http.get<any[]>(
      `${environment.apiUrl}/api/chart/area/${reference}`
    );
  }
  public newDataset(device: Device) {
    return this.http.post<Device>(`${environment.apiUrl}/api/datasets`, device);
  }
  public updateDataset(dataset: Dataset) {
    return this.http.put<Dataset>(
      `${environment.apiUrl}/api/datasets`,
      dataset
    );
  }
}
