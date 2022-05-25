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
    // looping through thresholdswitch items
    this.citem.items.forEach((item: thresholdswitch) => {
      this.items.push(this.addThresholdForm(item));
    });
  }
  get items() {
    return this.form.get("items") as FormArray;
  }
  addThresholdForm(item: thresholdswitch) {
    const form = this.formBuilder.group({
      item_description: [item.item_description],
      operation: [item.operation],
      datasource: [item.datasource],
      threshold_limit: [item.threshold_limit],
      on: [item.on],
    });
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
