# godem -SRTM HGT Reader and GDAL Interpolation Library

This is an advanced library designed to read Shuttle Radar Topography Mission (SRTM) HGT files and use the Geospatial Data Abstraction Library (GDAL) to perform elevation interpolation. This library aims to simplify the process of acquiring topographic data and translating it into practical, usable formats.

## Key Features

* **Caching of Downloads**: This feature reduces data redundancy and speeds up the reading process by storing previously downloaded data for future use. The next time you need the same HGT file, it retrieves the data from cache instead of initiating a new download.

* **Resolution Support**: The library is compatible with both 30m (1 arc-second) and 90m (3 arc-second) resolution files. This range of resolution support allows for versatility in your topographic data manipulation and visualization.

* **Data Source Support**: Our library can extract data from three primary sources:

    * [GPXSee](https://www.gpxsee.org/)
    * [ESA - European Space Agency](https://step.esa.int/main/)
    * [VIEWFINDER PANORAMAS](http://viewfinderpanoramas.org/)

These diverse sources ensure access to a comprehensive and accurate global topographic data set.

## Why Use This Library?

By incorporating GDAL for interpolation, our library brings together two powerful tools for processing and interpreting SRTM HGT files. It does so in a way that minimizes the technical barrier for users and maximizes the accessibility of the data.

Moreover, with caching functionality, the library also ensures efficient use of resources, saving both time and space during subsequent data processing.

## Getting Started

Detailed instructions on how to install and use the library will be provided soon. Stay tuned!

---

**NOTE**: This library is continuously updated and maintained. Contributions and suggestions for improvements are always welcome!
