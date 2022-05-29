import { Component, Input, OnInit } from "@angular/core";

@Component({
  selector: "app-two-column-div",
  templateUrl: "./two-column-div.component.html",
  styleUrls: ["./two-column-div.component.scss"],
})
export class TwoColumnDivComponent implements OnInit {
  @Input() col1: string | "no values";
  @Input() col2: string | "no values";

  constructor() {}

  ngOnInit(): void {}
}
