import { Component, Inject, OnInit } from "@angular/core";
import { MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA } from "@angular/material/legacy-dialog";

@Component({
  selector: "app-error-handling-dialog",
  templateUrl: "./error-handling-dialog.component.html",
  styleUrls: ["./error-handling-dialog.component.css"],
})
export class ErrorHandlingDialogComponent implements OnInit {
  constructor(@Inject(MAT_DIALOG_DATA) public data: any) {}

  ngOnInit(): void {}
}
