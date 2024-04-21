export function saveDraftTemplate(sts: number, edges: number): string {
  return (
    `Are you sure you want to save?: ` +
    `<ul>` +
    `<li>${sts} stations</li>` +
    `<li>${edges} edges</li>` +
    `</ul>`
  );
}
