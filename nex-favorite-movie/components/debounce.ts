export default function debounce(
  func: Function, 
  timeout = 300,
  //: not works with `setState`s
  withoutTimeoutActions: Function = () => { }) {
  let timer: NodeJS.Timeout;
  return (...args: any[]) => {
    clearTimeout(timer);
    timer = setTimeout(() => { func(...args); }, timeout);

    withoutTimeoutActions(...args);
  };
}