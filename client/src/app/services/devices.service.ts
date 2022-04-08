import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from 'src/environments/environment';
import { Device } from '../models/device';

@Injectable({
  providedIn: 'root',
})
export class DevicesService {
  constructor(private http: HttpClient) {}

  public LoadDevices() {
    return this.http.get<Device[]>(`${environment.apiUrl}/api/sensors`);
  }
  public LoadDevice(arduino_key: string) {
    return this.http.get<Device>(
      `${environment.apiUrl}/api/sensors/${arduino_key}`
    );
  }
  public AddNewDevice(title: string, description: string) {
    const obj = { title, description };
    return this.http.post<Device>(`${environment.apiUrl}/api/sensors`, obj);
  }
}
