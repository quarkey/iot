import { Component, Input, OnInit } from "@angular/core";
import { FormArray, FormBuilder, FormGroup, Validators } from "@angular/forms";
import { Controller } from "src/app/models/controllers";
import { Sensordata } from "src/app/models/dataset";
import { ControllersService } from "src/app/services/controllers.service";
import { environment } from "src/environments/environment";

@Component({
  selector: "app-controller-details",
  templateUrl: "./controller-details.component.html",
  styleUrls: ["./controller-details.component.scss"],
})
export class ControllerDetailsComponent implements OnInit {
  @Input() citem: Controller;
  form: FormGroup;
  categories: string[] = ["switch", "thresholdswitch", "timeswitch"];

  constructor(
    private formBuilder: FormBuilder,
    private controllerService: ControllersService
  ) {}
  defaultValue = { hour: 13, minute: 30 };

  timeChangeHandler(event: Event) {
    console.log(event);
  }

  invalidInputHandler() {
    // some error handling
  }
  ngOnInit(): void {
    this.form = this.formBuilder.group({
      category: [this.citem.category, Validators.required],
      title: [this.citem.title, Validators.required],
      description: [this.citem.description, Validators.required],
      items: this.formBuilder.array([]),
      active: [this.citem.active],
    });

    this.citem.items.forEach((item: any) => {
      switch (this.citem.category) {
        case "thresholdswitch":
          this.items.push(this.addThresholdSwitchForm(item as thresholdswitch));
          break;
        case "timeswitch":
          this.items.push(this.addTimeSwitchForm(item as timeswitch));
          break;
        case "switch":
          this.items.push(this.addSwitchForm(item as normalswitch));
          break;
      }
    });
  }
  get items() {
    return this.form.get("items") as FormArray;
  }
  addSwitchForm(item: any) {
    return this.formBuilder.group({
      item_description: [item.item_description || null, Validators.required],
      on: [item.on || null],
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
  updateController() {
    var obj = {
      ...this.form.value,
      id: this.citem.id,
    };

    this.controllerService.UpdateControllerByID(obj).subscribe((res) => {
      if (res) {
        this.form.markAsPristine();
      }
    });
  }
}

interface thresholdswitch {
  on: boolean;
  item_description: string;
  datasource: string;
  threshold_limit: number;
  operation: string;
}
interface timeswitch {
  on: boolean;
  item_description: string;
  time_on: string;
  time_off: string;
  duration: string;
  repeat: number;
}

interface normalswitch {
  on: boolean;
  item_description: string;
}
