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
        default:
          break;
      }
    });
  }
  get items() {
    return this.form.get("items") as FormArray;
  }
  addThresholdSwitchForm(item: thresholdswitch) {
    let form;
    if (item === null) {
      form = this.formBuilder.group({
        item_description: [null],
        operation: [null],
        datasource: [null],
        threshold_limit: [null],
        on: [null],
      });
    } else {
      form = this.formBuilder.group({
        item_description: [item.item_description],
        operation: [item.operation],
        datasource: [item.datasource],
        threshold_limit: [item.threshold_limit],
        on: [item.on],
      });
    }
    return form;
  }
  addTimeSwitchForm(item: timeswitch) {
    let form;
    if (item === null) {
      form = this.formBuilder.group({
        on: [null],
        repeat: [null],
        time_on: [null],
        time_off: [null],
        duration: [null],
        item_description: [null],
      });
    } else {
      form = this.formBuilder.group({
        on: [item.on],
        repeat: [item.repeat],
        time_on: [item.time_on],
        time_off: [item.time_off],
        duration: [item.duration],
        item_description: [item.item_description],
      });
    }
    return form;
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
