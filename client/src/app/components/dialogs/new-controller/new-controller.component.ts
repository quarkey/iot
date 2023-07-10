import { Component, OnInit } from "@angular/core";
import { FormBuilder, FormGroup, Validators } from "@angular/forms";
import { MatLegacyDialogRef as MatDialogRef } from "@angular/material/legacy-dialog";
import { Router } from "@angular/router";
import { ControllersService } from "src/app/services/controllers.service";

@Component({
  selector: "app-new-controller",
  templateUrl: "./new-controller.component.html",
  styleUrls: ["./new-controller.component.scss"],
})
export class NewControllerComponent implements OnInit {
  form: FormGroup;
  id: number;
  constructor(
    private formBuilder: FormBuilder,
    private dialogRef: MatDialogRef<NewControllerComponent>,
    private router: Router,
    private controllerService: ControllersService
  ) {}

  ngOnInit(): void {
    this.form = this.formBuilder.group({
      category: ["", Validators.required],
      title: ["", Validators.required],
      description: ["", Validators.required],
    });
  }
  addNewController() {
    var newDataset = { ...this.form.value };
    this.controllerService.newController(newDataset).subscribe((res) => {
      if (res) {
        console.log(res);
        this.router.navigate([`/controllers/${res}`]);
        this.dialogRef.close();
      }
    });
  }
}
