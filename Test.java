public class Test {
    public static int aint = 26;
    private static String bstr = ".. no more young";
    private static String convert(int vara, String varb) {
        return vara + varb;
    }

    public static void main(String[] args) {
        String res = convert(aint, bstr);
        System.out.println(res);
    }
}
