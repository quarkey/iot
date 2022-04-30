import { Injectable } from "@angular/core";
import { environment } from "src/environments/environment";

@Injectable({
  providedIn: "root",
})
export class ToolbarService {
  pageTitle = environment.title;
  constructor() {}

  setTitle(title: string) {
    this.pageTitle = title;
  }
  getTitle() {
    return this.pageTitle;
  }
}
