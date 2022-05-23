import { Component, Input, OnInit } from "@angular/core";
import { FormArray, FormBuilder, FormGroup, Validators } from "@angular/forms";
import { Controller } from "src/app/models/controllers";
import { ControllersService } from "src/app/services/controllers.service";

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
    const thresholdForm = this.formBuilder.group({
      on: [this.citem.items[0].on],
      item_description: [
        this.citem.items[0].item_description,
        Validators.required,
      ],
      operation: [this.citem.items[0].operation, Validators.required],
      datasource: [this.citem.items[0].datasource, Validators.required],
      threshold_limit: [
        this.citem.items[0].threshold_limit,
        Validators.required,
      ],
    });

    this.form = this.formBuilder.group({
      category: [this.citem.category, Validators.required],
      title: [this.citem.title, Validators.required],
      description: [this.citem.description, Validators.required],
      items: thresholdForm,
      active: [this.citem.active],
    });
  }
  updateController() {
    var obj = {
      ...this.form.value,
      id: this.citem.id,
    };
    var temp = [obj.items];
    obj.items = temp;

    this.controllerService.UpdateControllerByID(obj).subscribe((res) => {
      if (res) {
        this.form.markAsPristine();
      }
    });
  }
}
