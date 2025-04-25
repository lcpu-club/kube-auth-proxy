export interface Partition {
  Name: string;
  Description: string;
  GPUTag: string;
  GPUName: string;
  Images: string[];
  CPULimit: number;
  MemoryLimit: number; // in GiB
}
