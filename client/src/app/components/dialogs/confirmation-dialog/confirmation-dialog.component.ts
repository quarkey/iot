import { Component, Inject, OnInit } from "@angular/core";
import { FormControl, Validators } from "@angular/forms";
import { MatDialogRef, MAT_DIALOG_DATA } from "@angular/material/dialog";

@Component({
  selector: "app-confirmation-dialog",
  templateUrl: "./confirmation-dialog.component.html",
  styleUrls: ["./confirmation-dialog.component.css"],
})
export class ConfirmationDialogComponent implements OnInit {
  constructor(
    public dialogRef: MatDialogRef<ConfirmationDialogComponent>,
    @Inject(MAT_DIALOG_DATA)
    public data: {
      title: string;
      message: string;
      yes: string;
      cancel: string;
      confirmForm: boolean;
    }
  ) {}
  form: FormControl;
  ngOnInit(): void {
    if (this.data.confirmForm) {
      this.form = new FormControl("", [
        Validators.required,
        Validators.pattern("confirm"),
      ]);
    }
  }
  click(clickmessage: string): void {
    this.dialogRef.close(clickmessage);
  }
}
