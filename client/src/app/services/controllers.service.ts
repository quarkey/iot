import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { FormBuilder, Validators } from "@angular/forms";
import { environment } from "src/environments/environment";
import { Controller, thresholdswitch, timeswitch } from "../models/controllers";

@Injectable({
  providedIn: "root",
})
export class ControllersService {
  constructor(private http: HttpClient, private formBuilder: FormBuilder) {}

  public LoadControllersList() {
    return this.http.get<Controller[]>(`${environment.apiUrl}/api/controllers`);
  }
  public LoadControllerByID(id: number) {
    return this.http.get<Controller>(
      `${environment.apiUrl}/api/controllers/${id}`
    );
  }
  public UpdateControllerByID(citem: Controller) {
    return this.http.put<Controller>(
      `${environment.apiUrl}/api/controllers`,
      citem
    );
  }
  public DeleteControllerBy(id: number) {
    return this.http.delete<any>(
      `${environment.apiUrl}/api/controllers/${id}`,
      null
    );
  }
  public LoadDefualtControllerItemsValues(id: number) {
    return this.http.post<any>(
      `${environment.apiUrl}/api/controllers/${id}/reset`,
      null
    );
  }
  public newController(citem: any) {
    return this.http.post<any>(`${environment.apiUrl}/api/controllers`, citem);
  }

  public setContllerState(id: number, switchState: string) {
    return this.http.get<any>(
      `${environment.apiUrl}/api/controller/${id}/switch/${switchState}`
    );
  }
  addInitialForm(item: Controller) {
    return this.formBuilder.group({
      category: [item.category, Validators.required],
      title: [item.title, Validators.required],
      description: [item.description, Validators.required],
      items: this.formBuilder.array([]),
      active: [item.active],
    });
  }
  addSwitchForm(item: any) {
    return this.formBuilder.group({
      item_description: [item.item_description || null, Validators.required],
      on: [item.on || null],
    });
  }
  addTimeSwitchForm(item: timeswitch) {
    return this.formBuilder.group({
      on: [item.on || null],
      repeat: [item.repeat || null, Validators.required],
      time_on: [item.time_on || null, Validators.required],
      time_off: [item.time_off || null, Validators.required],
      duration: [item.duration || null, Validators.required],
      item_description: [item.item_description || null, Validators.required],
    });
  }
  addThresholdSwitchForm(item: thresholdswitch) {
    return this.formBuilder.group({
      item_description: [item.item_description || null, Validators.required],
      operation: [item.operation || null, Validators.required],
      datasource: [item.datasource || null, Validators.required],
      threshold_limit: [item.threshold_limit || null, Validators.required],
      on: [item.on || null],
    });
  }
}
