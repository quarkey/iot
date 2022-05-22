import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { environment } from "src/environments/environment";
import { Controller } from "../models/controllers";

@Injectable({
  providedIn: "root",
})
export class ControllersService {
  constructor(private http: HttpClient) {}

  public LoadControllersList() {
    return this.http.get<Controller[]>(`${environment.apiUrl}/api/controllers`);
  }
  public LoadControllerByID(id: number) {
    return this.http.get<Controller>(
      `${environment.apiUrl}/api/controllers/${id}`
    );
  }
}
