import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from 'src/environments/environment';
import { Dataset } from '../models/dataset';

@Injectable({
  providedIn: 'root',
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
    return this.http.get<any[]>(
      `${environment.apiUrl}/api/sensordata/${reference}`
    );
  }
}
