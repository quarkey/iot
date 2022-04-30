import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { environment } from "src/environments/environment";
import { Dashboard } from "../models/general";

@Injectable({
  providedIn: "root",
})
export class GeneralService {
  constructor(private http: HttpClient) {}

  public Dashboard() {
    return this.http.get<Dashboard>(`${environment.apiUrl}/api/dashboard`);
  }
}
