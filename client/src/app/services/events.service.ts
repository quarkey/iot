import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { environment } from "src/environments/environment";
import { EventList } from "../models/event";

@Injectable({
  providedIn: "root",
})
export class EventsService {
  constructor(private http: HttpClient) {}

  public LoadEvents(count: number) {
    return this.http.get<EventList[]>(
      `${environment.apiUrl}/api/events/${count}`
    );
  }
}
