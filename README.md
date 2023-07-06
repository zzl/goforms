# GoForms

## What is this?
GoForms is a GUI programming framework designed for use with the Go language.
It provides similar functionality to WinForms.

- It simplifies window message handling.
- It includes object-oriented wrappers for common controls.
- It features a layout system.
- It includes a drawing package based on the gdiplus technology.

## Design Principals

- **Lightweight Functions:** UI functions should be simple wrappers
  over the native parts of the system.
- **Lightweight Objects:** UI objects should not maintain unnecessary state data if 
  they are already contained in the underlying native components.
- **Easy to use API:** The API should be designed to accommodate the syntax of Go,  
  making function calls smooth and concise.   
- **Easy to extend:** The code should follow an object-oriented style, 
  allowing for component inheritance and overriding.
  > While OOP may not be widely embraced, particularly within the Go community, 
  > I think it offers a suitable approach for this framework.

## Screenshots

* Screenshot of the fileexplore example
![fileexplore.png](https://github.com/zzl/goforms/blob/assets/fileexplorer.png?raw=true)

* Screen recording of the goforms-designer app (still in early development)
![designer.gif](https://github.com/zzl/goforms/blob/assets/designer.gif?raw=true)

